package octoprint

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

// Octoprint - plugins main structure
type Octoprint struct {
	URL    string `toml:"url"`
	APIKey string `toml:"apikey"`
	API    OctoAPI
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
`
}

// Init setup the octoprint client
func (o *Octoprint) Init() error {
	o.API = NewGoOcto(o.URL, o.APIKey)
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

	return nil
}

func init() {
	inputs.Add("octoprint", func() telegraf.Input { return &Octoprint{} })
}
