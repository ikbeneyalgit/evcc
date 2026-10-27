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
}
