package octoprint

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	octoapi "github.com/mcuadros/go-octoprint"
)

// Octoprint - plugins main structure
type Octoprint struct {
	URL    string `toml:"url"`
	APIKey string `toml:"apikey"`
	client *octoapi.Client
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
	o.client = octoapi.NewClient(o.URL, o.APIKey)
	return nil
}

// Gather OctoPrint metrics
func (o *Octoprint) Gather(acc telegraf.Accumulator) error {

	state, _ := o.getState()

	acc.AddFields("state",
		map[string]interface{}{
			"value": state,
		},
		map[string]string{
			"id": "State",
		},
	)

	tools, _ := o.getToolInfo()

	for _, t := range tools {
		acc.AddFields("tool",
			map[string]interface{}{
				"name":       t.name,
				"actualTemp": t.actualTemp,
				"targetTemp": t.targetTemp,
			},
			map[string]string{
				"name": t.name,
			},
		)
	}

	return nil
}

func (o *Octoprint) getState() (string, error) {
	r := octoapi.ConnectionRequest{}
	s, err := r.Do(o.client)
	if err != nil {
		return "", err
	}

	return string(s.Current.State), nil
}

type tool struct {
	name       string
	actualTemp float64
	targetTemp float64
}

func (o *Octoprint) getToolInfo() ([]tool, error) {
	r := octoapi.StateRequest{}
	s, err := r.Do(o.client)
	if err != nil {
		return []tool{}, err
	}

	var tools []tool
	for toolName, state := range s.Temperature.Current {
		var t tool
		t.name = toolName
		t.actualTemp = state.Actual
		t.targetTemp = state.Target
		tools = append(tools, t)
	}

	return tools, nil
}

func init() {
	inputs.Add("octoprint", func() telegraf.Input { return &Octoprint{} })
}
