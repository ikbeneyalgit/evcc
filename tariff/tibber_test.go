package tariff

import (
	"testing"
	"time"

	"github.com/evcc-io/evcc/meter/tibber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTibberChargesZones guards that charge zones alone switch to the energy price base (#34558)
func TestTibberChargesZones(t *testing.T) {
	ts := time.Date(2026, 1, 15, 2, 0, 0, 0, time.Local)
	prices := []tibber.Price{{StartsAt: ts, Total: 0.30, Energy: 0.20}}

	tb := &Tibber{embed: &embed{}}
	require.NoError(t, tb.init())
	assert.InDelta(t, 0.30, tb.rates(prices)[0].Value, 1e-9, "total without charges")

	tb = &Tibber{embed: &embed{
		ChargesZones_: []chargesZoneConfig{{Charges: -0.07, Hours: "00:00-03:00"}},
	}}
	require.NoError(t, tb.init())
	assert.InDelta(t, 0.20-0.07, tb.rates(prices)[0].Value, 1e-9, "energy with zone charges")
}

// TestTibberPollInterval covers polling more often from noon until tomorrow's prices arrive
func TestTibberPollInterval(t *testing.T) {
	at := func(h, m int) time.Time {
		return time.Date(2026, 10, 28, h, m, 0, 0, time.Local)
	}

	for _, tc := range []struct {
		now      time.Time
		tomorrow bool
		wait     time.Duration
	}{
		{at(11, 59), false, time.Hour},
		{at(12, 0), false, 5 * time.Minute},
		{at(13, 30), false, 5 * time.Minute},
		{at(23, 55), false, 5 * time.Minute},
		{at(13, 30), true, time.Hour},
		{at(9, 0), true, time.Hour},
	} {
		assert.Equal(t, tc.wait, tibberPollInterval(tc.now, tc.tomorrow), "%s tomorrow=%v", tc.now.Format("15:04"), tc.tomorrow)
	}
}
