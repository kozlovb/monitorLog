package alert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func registerTimestampTimes(timestamp int, times int, a *Alert) {
	for i := 0; i < times; i++ {
		a.RegisterEntry(timestamp)
	}
}

func checkAlert(t *testing.T, a *Alert, expected_state bool) {
	actual_state := a.GetAlertState()
	require.Equal(t, expected_state, actual_state, "Expected alert state %t, actual %t\n", expected_state, actual_state)
}

func checkAverageEquals(t *testing.T, a *Alert, expected_average float64) {
	actual_req_per_econd := a.GetAverageRequestPerSecond()
	require.InDelta(t, expected_average, actual_req_per_econd, 0.01, "Expected average request per second: %.2f is not within 0.01 margin of the actual one: %.2f\n", expected_average, actual_req_per_econd)
}

func Test_alerState(t *testing.T) {

	timeInterval := 10
	threshold := 1

	alert := NewAlert(timeInterval, threshold)

	alert.RegisterEntry(2)
	registerTimestampTimes(4, 2, alert)

	checkAlert(t, alert, false)
	checkAverageEquals(t, alert, 0.3)

	registerTimestampTimes(10, 8, alert)

	checkAlert(t, alert, true)
	checkAverageEquals(t, alert, 1.1)

	registerTimestampTimes(18, 3, alert)

	checkAlert(t, alert, true)
	checkAverageEquals(t, alert, 1.1)

	alert.RegisterEntry(50)

	checkAlert(t, alert, false)
	checkAverageEquals(t, alert, 0.1)
}
