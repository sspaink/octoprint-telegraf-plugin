package octoprint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq" // I assume we don't want this in the main.go shim
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
	DB         *sqlx.DB
	Log        telegraf.Logger
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
		DB, err := sqlx.Open("postgres", URI)
		if err != nil {
			o.Log.Errorf("Failed to open the DB connection: %v", err)
			return err
		}
		o.DB = DB
	}

	return nil
}

// Gather OctoPrint metrics
func (o *Octoprint) Gather(acc telegraf.Accumulator) error {

	l, err := o.GetLayerProgress()
	if err == nil {
		o.UploadLayerProgress(l, acc)
	}

	state, err := o.State()
	if err == nil {
		o.UploadState(state, acc)
	}

	tools, err := o.ToolInfo()
	if err == nil {
		o.UploadToolInfo(tools, acc)
	}

	if o.DB != nil {
		o.GatherFilamentManagerData(acc)
	}

	return nil
}

// State returns the current state of the 3d printer (e.g. printing, operational, closed)
func (o *Octoprint) State() (string, error) {
	c, err := o.API.ConnectionRequest()
	if err != nil {
		o.Log.Errorf("Failed to get state information: %v", err)
		return "", err
	}
	return string(c.Current.State), nil
}

// UploadState will upload the state information to InfluxDB
func (o *Octoprint) UploadState(state string, acc telegraf.Accumulator) {
	acc.AddFields("state",
		map[string]interface{}{
			"value": state,
		},
		map[string]string{
			"id": "State",
		},
	)
}

// Tool defines a tool on the 3d printer (e.g. hotend)
type Tool struct {
	Name       string
	ActualTemp float64
	TargetTemp float64
}

// ToolInfo returns a list of tools on the connected 3d printer
func (o *Octoprint) ToolInfo() ([]Tool, error) {
	var tools []Tool
	s, err := o.API.StateRequest()
	if err != nil {
		o.Log.Errorf("Failed to make a state request to Octoprint: %v", err)
		return tools, err
	}

	for toolName, state := range s.Temperature.Current {
		var t Tool
		t.Name = toolName
		t.ActualTemp = state.Actual
		t.TargetTemp = state.Target
		tools = append(tools, t)
	}

	return tools, nil
}

// UploadToolInfo will upload tool information to InfluxDB
func (o *Octoprint) UploadToolInfo(tools []Tool, acc telegraf.Accumulator) {
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
}

func init() {
	inputs.Add("octoprint", func() telegraf.Input { return &Octoprint{} })
}

type LayerProgress struct {
	Layer LayerData `json:"layer"`
}

type LayerData struct {
	Current string `json:"current"`
	Total   string `json:"total"`
}

func (o *Octoprint) GetLayerProgress() (LayerProgress, error) {
	URL := o.URL + "/plugin/DisplayLayerProgress/values"
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		o.Log.Errorf("Failed creating new request for layer progress: %v", err)
		return LayerProgress{}, err
	}
	req.Header.Set("Authorization", "Bearer "+o.APIKey)
	res, err := client.Do(req)
	if err != nil {
		o.Log.Errorf("Failed to execute request for layer progress: %v", err)
		return LayerProgress{}, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		o.Log.Errorf("Failed to read response body for layer progress: %v", err)
		return LayerProgress{}, err
	}
	var l LayerProgress
	err = json.Unmarshal(body, &l)
	if err != nil {
		o.Log.Errorf("Failed to unmarshall layer progress response: %v", err)
		return LayerProgress{}, err
	}
	return l, nil
}

func (o *Octoprint) UploadLayerProgress(l LayerProgress, acc telegraf.Accumulator) {
	var current, total int

	if l.Layer.Current == "-" {
		current = 0
	} else {
		var err error
		current, err = strconv.Atoi(l.Layer.Current)
		if err != nil {
			o.Log.Errorf("Failed to convert current layer to int: %v", err)
			return
		}
	}

	if l.Layer.Total == "-" {
		total = 0
	} else {
		var err error
		total, err = strconv.Atoi(l.Layer.Total)
		if err != nil {
			o.Log.Errorf("Failed to convert total layer to int: %v", err)
		}
	}

	acc.AddFields("layer progress",
		map[string]interface{}{
			"current_layer": current,
			"total_layer":   total,
		},
		map[string]string{
			"layers": "active layers",
		},
	)
}
