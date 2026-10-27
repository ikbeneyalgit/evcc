package meter

import (
	"testing"

	"github.com/evcc-io/evcc/api"
	"github.com/stretchr/testify/require"
)

// TestVictronEnergyBatteryModes covers the battery modes announced by the template
func TestVictronEnergyBatteryModes(t *testing.T) {
	m, err := NewFromConfig(t.Context(), "template", map[string]any{
		"template": "victron-energy",
		"usage":    "battery",
		"modbus":   "tcpip",
		"host":     "localhost",
		"port":     502,
		"minsoc":   20,
	})
	require.NoError(t, err)

	ctrl, ok := api.Cap[api.BatteryController](m)
	require.True(t, ok)
	require.Equal(t, []api.BatteryMode{api.BatteryNormal, api.BatteryHold, api.BatteryCharge}, ctrl.BatteryModes())

	// maxdischargepower does not lead to BatteryDischarge (for users who want to keep Victron DESS in charge)
	m, err = NewFromConfig(t.Context(), "template", map[string]any{
		"template":          "victron-energy",
		"usage":             "battery",
		"modbus":            "tcpip",
		"host":              "localhost",
		"port":              502,
		"minsoc":            20,
		"maxdischargepower": 3000,
	})
	require.NoError(t, err)

	ctrl, ok = api.Cap[api.BatteryController](m)
	require.True(t, ok)
	require.Equal(t, []api.BatteryMode{api.BatteryNormal, api.BatteryHold, api.BatteryCharge}, ctrl.BatteryModes())

	// BatteryDischarge supported if user opted into letting evcc manage grid discharge (on this battery)
	m, err = NewFromConfig(t.Context(), "template", map[string]any{
		"template":          "victron-energy",
		"usage":             "battery",
		"modbus":            "tcpip",
		"host":              "localhost",
		"port":              502,
		"minsoc":            20,
		"maxdischargepower": 3000,
		"griddischarge":     true,
	})
	require.NoError(t, err)

	ctrl, ok = api.Cap[api.BatteryController](m)
	require.True(t, ok)
	require.Equal(t, []api.BatteryMode{api.BatteryNormal, api.BatteryHold, api.BatteryCharge, api.BatteryDischarge}, ctrl.BatteryModes())

	// griddischarge without maxdischargepower discharges at full power
	m, err = NewFromConfig(t.Context(), "template", map[string]any{
		"template":      "victron-energy",
		"usage":         "battery",
		"modbus":        "tcpip",
		"host":          "localhost",
		"port":          502,
		"minsoc":        20,
		"griddischarge": true,
	})
	require.NoError(t, err)

	ctrl, ok = api.Cap[api.BatteryController](m)
	require.True(t, ok)
	require.Equal(t, []api.BatteryMode{api.BatteryNormal, api.BatteryHold, api.BatteryCharge, api.BatteryDischarge}, ctrl.BatteryModes())
}
