package octoprint

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"

	_ "github.com/lib/pq"
)

// Octoprint - plugins main structure
type Octoprint struct {
	URL        string `toml:"url"`
	APIKey     string `toml:"apikey"`
	DBNamePSQL string `toml:"dbnamepsql"`
	UserPSQL   string `toml:"userpsql"`
	PassPSQL   string `toml:"passpsql"`
	IP         string `tom:"ip"`
	API        OctoAPI
	DB         *sql.DB
}

// Description returns the plugin description
func (o *Octoprint) Description() string {
	return "A plugin to gather data from OctoPrint"
}

// SampleConfig returns sample configuration for this plugin
func (o *Octoprint) SampleConfig() string {
	return `
  ## Indicate if everything is fine
  [inputs.octoprint]
  ## OctoPrint's URL
  # url=""
  ## OctoPrint's API Key
  # apikey=""
  ## OPTIONAL IF Filamanet Manager PSQL IS SET
  ## PSQL Database name
  # dbnamepsql=""
  ## Username that has access to the database
  # userpsql=""
  ## Password for the user
  # passpsql=""
  ## IP that is hosting the database
  # ip=""
`
}

func (o *Octoprint) verifyPSQLSettings() bool {
	if o.DBNamePSQL != "" && o.UserPSQL != "" && o.PassPSQL != "" && o.IP != "" {
		return true
	}
	return false
}

// Init setup the octoprint client
func (o *Octoprint) Init() error {
	o.API = NewGoOcto(o.URL, o.APIKey)
	// If PSQL is set in the settings, create a connection
	if o.verifyPSQLSettings() {
		URI := fmt.Sprintf("postgres://%s:%s@%s/%s", o.UserPSQL, o.PassPSQL, o.IP, o.DBNamePSQL)
		DB, err := sql.Open("postgres", URI)
		if err != nil {
			log.Fatal("Failed to open a DB connection: ", err)
		}
		o.DB = DB
	}

	return nil
}

// Tool defines a tool on the 3d printer (e.g. hotend)
type Tool struct {
	Name       string
	ActualTemp float64
	TargetTemp float64
}

// ToolInfo returns a list of tools on the connected 3d printer
func (o *Octoprint) ToolInfo() ([]Tool, error) {
	s, err := o.API.StateRequest()
	if err != nil {
		return []Tool{}, err
	}

	var tools []Tool
	for toolName, state := range s.Temperature.Current {
		var t Tool
		t.Name = toolName
		t.ActualTemp = state.Actual
		t.TargetTemp = state.Target
		tools = append(tools, t)
	}

	return tools, nil
}

// State returns the current state of the 3d printer (e.g. printing)
func (o *Octoprint) State() (string, error) {
	c, err := o.API.ConnectionRequest()
	if err != nil {
		return "", err
	}

	return string(c.Current.State), nil
}

// Gather OctoPrint metrics
func (o *Octoprint) Gather(acc telegraf.Accumulator) error {

	state, _ := o.State()

	acc.AddFields("state",
		map[string]interface{}{
			"value": state,
		},
		map[string]string{
			"id": "State",
		},
	)

	tools, _ := o.ToolInfo()

	for _, t := range tools {
		acc.AddFields("tool",
			map[string]interface{}{
				"name":       t.Name,
				"actualTemp": t.ActualTemp,
				"targetTemp": t.TargetTemp,
			},
			map[string]string{
				"name": t.Name,
			},
		)
	}

	// If there is a postgress connection
	if o.DB != nil {
		rows, _ := o.DB.Query("SELECT vendor, material FROM profiles;")

		for rows.Next() {
			var vendor string
			var material string
			rows.Scan(&vendor, &material)
			acc.AddFields("filament",
				map[string]interface{}{
					"vendor":   vendor,
					"material": material,
				},
				map[string]string{
					"id": "FilamentManagerProfiles",
				},
			)
		}
	}

	return nil
}

func init() {
	inputs.Add("octoprint", func() telegraf.Input { return &Octoprint{} })
}
