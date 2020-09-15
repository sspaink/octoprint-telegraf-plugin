package octoprint

import (
	"fmt"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	octoapi "github.com/mcuadros/go-octoprint"
	"github.com/prometheus/common/log"
)

// Octoprint - plugins main structure
type Octoprint struct {
	URL    string `toml:"url"`
	APIKey string `toml:"apikey"`
	client *octoapi.Client
}

func (o *Octoprint) Description() string {
	return "A plugin to gather data from OctoPrint"
}

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

func (o *Octoprint) Init() error {
	return nil
}

func (o *Octoprint) Gather(acc telegraf.Accumulator) error {
	if o.client == nil {
		o.client = octoapi.NewClient(o.URL, o.APIKey)
	}

	r := octoapi.ConnectionRequest{}
	s, err := r.Do(o.client)
	if err != nil {
		log.Error("error requesting connection state: %s", err)
	}

	fmt.Printf("Connection State: %q\n", s.Current.State)

	acc.AddFields("state", map[string]interface{}{"value": s.Current.State}, nil)

	return nil
}

func init() {
	inputs.Add("octoprint", func() telegraf.Input { return &Octoprint{} })
}
