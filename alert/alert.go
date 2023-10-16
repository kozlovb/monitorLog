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

func (a *Alert) RegisterEntry(timestamp int) {
	a.findUpperBoundAndAdd(timestamp)
	a.resetOldEntries(timestamp)
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

// vector with rotaing indexes
// for sure there hopuld be somw way to  look for an element in a sorted lit fater than 0(n)
// non nil is required
func (a *Alert) findUpperBoundAndAdd(timestamp int) {

	//check if container is empty
	e := a.container_for_logs_sliding_window.Front()
	if e == nil {
		a.container_for_logs_sliding_window.PushBack(AlertEntry{time_in_seconds: timestamp, qty: 1})
		a.total_hits++
		return
	}
	e_value, _ := e.Value.(AlertEntry)

	//check if earlier than first timestamp
	if timestamp < e_value.time_in_seconds {
		return
	}

	//check if later or equal to the last timestamp
	// as it's the most common case
	last_el := a.container_for_logs_sliding_window.Back()
	last_value, _ := last_el.Value.(AlertEntry)
	if timestamp > last_value.time_in_seconds {
		a.container_for_logs_sliding_window.PushBack(AlertEntry{time_in_seconds: timestamp, qty: 1})
		a.total_hits++
		return
	} else if timestamp == last_value.time_in_seconds {
		last_value.qty++
		last_el.Value = last_value
		a.total_hits++
		return
	}

	//should be with in the bounds
	for e != nil {
		e_value = e.Value.(AlertEntry)
		if e_value.time_in_seconds < timestamp {
			e = e.Next()
		} else {
			break
		}
	}

	if e_value.time_in_seconds == timestamp {
		e_value.qty++
		e.Value = e_value
	} else {
		a.container_for_logs_sliding_window.InsertBefore(AlertEntry{time_in_seconds: timestamp, qty: 1}, e)
	}
	a.total_hits++
}

func (a *Alert) resetOldEntries(timestamp int) {
	for e := a.container_for_logs_sliding_window.Front(); e != nil; {
		next := e.Next()
		e_value := e.Value.(AlertEntry)
		if e_value.time_in_seconds+a.time_interval_for_average <= timestamp {
			a.total_hits = a.total_hits - e_value.qty
			a.container_for_logs_sliding_window.Remove(e)
		} else {
			break
		}
		e = next
	}
}

// 2019-02-07T21:12:36Z  - first alert
//2019-02-07T21:16:03Z  - second alert
