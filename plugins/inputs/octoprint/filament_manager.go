package octoprint

import (
	"fmt"

	"github.com/influxdata/telegraf"
)

// GatherFilamentManagerData will gather/upload data stored by the filamant manager extension
// The filament manager is a useful extension to track your filament usage
// This extension supports uploading the data to an external Postgress database, which this plugin can read from
// More details: https://plugins.octoprint.org/plugins/filamentmanager/
func (o *Octoprint) GatherFilamentManagerData(acc telegraf.Accumulator) {
	spool, err := o.SelectedSpool()
	if err == nil {
		o.UploadSpoolData(spool, acc)
	}
}

// SelectedSpool defines a spool of filament currently selected to be used
type SelectedSpool struct {
	ID       string
	Name     string
	Weight   int
	Used     int
	Vendor   string
	Material string
}

// SelectedSpool will gather the data for the currently selected spool
func (o *Octoprint) SelectedSpool() (SelectedSpool, error) {
	var spool SelectedSpool
	query := `
	SELECT spools.id, spools.name, spools.weight, spools.used, profiles.vendor, profiles.material 
	FROM spools, profiles, selections 
	WHERE spools.profile_id = profiles.id AND selections.spool_id = spools.id;
	`
	row := o.DB.QueryRowx(query)
	err := row.StructScan(&spool)
	if err != nil {
		o.Log.Errorf("Failed to query for current spool information: %v", err)
		return spool, err
	}

	return spool, nil
}

// UploadSpoolData will upload the information for the currently selected spool to InfluxDB
func (o *Octoprint) UploadSpoolData(spool SelectedSpool, acc telegraf.Accumulator) {
	acc.AddFields("filament",
		map[string]interface{}{
			"name": fmt.Sprintf("Material: %s Color: %s", spool.Material, spool.Name),
			"used": spool.Used,
		},
		map[string]string{
			"id": fmt.Sprintf("%s_%s", spool.ID, spool.Name),
		},
	)
}
