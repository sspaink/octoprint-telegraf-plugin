package octoprint

import (
	"testing"

	gooctoprint "github.com/mcuadros/go-octoprint"
)

type MockOctoAPI struct{}

func (m *MockOctoAPI) StateRequest() (*gooctoprint.FullStateResponse, error) {
	var t gooctoprint.TemperatureData
	t.Actual = 200
	t.Target = 200

	var r gooctoprint.FullStateResponse
	r.Temperature.Current = make(map[string]gooctoprint.TemperatureData)
	r.Temperature.Current["tool0"] = t

	return &r, nil
}

func (m *MockOctoAPI) ConnectionRequest() (*gooctoprint.ConnectionResponse, error) {
	var c gooctoprint.ConnectionResponse
	c.Current.State = "printing"
	return &c, nil
}

func TestToolInfo(t *testing.T) {
	o := Octoprint{}
	o.API = &MockOctoAPI{}
	toolInfo, err := o.ToolInfo()
	if err != nil {
		t.Errorf("Failed to get tool info, %s", err)
	}

	if len(toolInfo) == 0 {
		t.Errorf("No tools returned")
	}

	var expectedTool Tool
	expectedTool.ActualTemp = 200
	expectedTool.TargetTemp = 200
	expectedToolName := "tool0"

	if toolInfo[0].ActualTemp != expectedTool.ActualTemp {
		t.Errorf("Wrong actual temp returned, got: %f, want: %f", toolInfo[0].ActualTemp, expectedTool.ActualTemp)
	}
	if toolInfo[0].TargetTemp != expectedTool.TargetTemp {
		t.Errorf("Wrong actual temp returned, got: %f, want: %f", toolInfo[0].TargetTemp, expectedTool.TargetTemp)
	}
	if toolInfo[0].Name != expectedToolName {
		t.Errorf("Wrong actual temp returned, got: %s, want: %s", toolInfo[0].Name, expectedToolName)
	}

}

func TestState(t *testing.T) {
	o := Octoprint{}
	o.API = &MockOctoAPI{}
	state, err := o.State()
	if err != nil {
		t.Errorf("Failed to get state, %s", err.Error())
	}

	expectedState := "printing"

	if state != expectedState {
		t.Errorf("Returned invalid state, got: %s, want: %s", state, expectedState)
	}
}
