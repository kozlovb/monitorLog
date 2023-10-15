package alert

import (
	"container/list"
	"fmt"
	"time"
)

type AlertEntry struct {
	time_in_seconds int
	qty             int
}

// TODO change to durationtime interval ?
type Alert struct {
	time_interval_for_average         int
	threshold_requests_per_second     int
	total_hits                        int
	container_for_logs_sliding_window *list.List
}

func NewAlert(timeInterval int, threshold int) *Alert {
	slidingWindow := list.New()
	return &Alert{
		time_interval_for_average:         timeInterval,
		threshold_requests_per_second:     threshold,
		container_for_logs_sliding_window: slidingWindow,
		total_hits:                        0,
	}
}

func (a *Alert) RegisterEntry(timestamp int) bool {

	lastElement := a.container_for_logs_sliding_window.Back()
	var last_value AlertEntry
	var ok bool

	if lastElement != nil {
		last_value, ok = lastElement.Value.(AlertEntry)
	}

	// here if log
	if !ok || last_value.time_in_seconds != timestamp {
		a.container_for_logs_sliding_window.PushBack(AlertEntry{time_in_seconds: timestamp, qty: 1})
	} else {
		last_value.qty++
		lastElement.Value = last_value
	}

	a.total_hits++
	a.resetOldEntries(timestamp)
	return a.container_for_logs_sliding_window.Len() > a.threshold_requests_per_second*a.time_interval_for_average
}

func (a *Alert) GetAlertState() bool {
	return a.total_hits > a.threshold_requests_per_second*a.time_interval_for_average
}

func (a *Alert) GetAverageRequestPerSecond() float64 {
	return float64(a.total_hits) / float64(a.time_interval_for_average)
}

func (a *Alert) GenerateAlertMsg(timestamp int) string {
	triggeredTime := time.Unix(int64(timestamp), 0).UTC()
	alertMessage := fmt.Sprintf("High traffic generated an alert - hits = %.2f, triggered at %s", a.GetAverageRequestPerSecond(), triggeredTime.Format(time.RFC3339))
	return alertMessage
}

func (a *Alert) GenerateRecoveryAlertMsg(timestamp int) string {
	triggeredTime := time.Unix(int64(timestamp), 0).UTC()
	alertMessage := fmt.Sprintf("Traffic has recovered - hits = %.2f, recovered at %s", a.GetAverageRequestPerSecond(), triggeredTime.Format(time.RFC3339))
	return alertMessage
}

func (a *Alert) resetOldEntries(timestamp int) {
	for e := a.container_for_logs_sliding_window.Front(); e != nil; {
		next := e.Next()
		e_value := e.Value.(AlertEntry)
		if e_value.time_in_seconds+a.time_interval_for_average < timestamp {
			a.total_hits = a.total_hits - e_value.qty
			a.container_for_logs_sliding_window.Remove(e)
		} else {
			break
		}
		e = next
	}
}
