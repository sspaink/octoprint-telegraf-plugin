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

	r := octoapi.ConnectionRequest{}
	s, err := r.Do(o.client)
	if err != nil {
		return err
	}

	acc.AddFields("state",
		map[string]interface{}{"value": string(s.Current.State)},
		map[string]string{"id": "State"},
	)

	return nil
}

func init() {
	inputs.Add("octoprint", func() telegraf.Input { return &Octoprint{} })
}
