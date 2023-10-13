package alarm

import (
	"container/list"
	"fmt"
	"time"
)

type AlarmEntry struct {
	time_in_seconds int
	qty             int
}

// TODO change to durationtime interval ?
type Alarm struct {
	time_interval_for_average         int
	threshold_requests_per_second     int
	total_hits                        int
	container_for_logs_sliding_window *list.List
}

func NewAlarm(timeInterval int, threshold int) *Alarm {
	// Create a new list for the sliding window
	slidingWindow := list.New()

	return &Alarm{
		time_interval_for_average:         timeInterval,
		threshold_requests_per_second:     threshold,
		container_for_logs_sliding_window: slidingWindow,
		total_hits:                        0,
	}
}

// TODO nepoporadku
// make NewAlarm
func (a *Alarm) RegisterEntry(timestamp int) bool {
	// check
	lastElement := a.container_for_logs_sliding_window.Back()
	//TODO one else as it does the same ?
	if lastElement != nil {
		last_value, ok := lastElement.Value.(AlarmEntry)
		//todo treat not ok
		if ok && last_value.time_in_seconds == timestamp {
			last_value.qty += 1
			lastElement.Value = last_value
		} else {
			a.container_for_logs_sliding_window.PushBack(AlarmEntry{time_in_seconds: timestamp, qty: 1})
		}
	} else {
		a.container_for_logs_sliding_window.PushBack(AlarmEntry{time_in_seconds: timestamp, qty: 1})
	}
	a.total_hits++
	// remove timestamp time_interval_for_average
	a.resetOldEntries(timestamp)
	if a.container_for_logs_sliding_window.Len() > a.threshold_requests_per_second*a.time_interval_for_average {
		//fmt.Println("total ", a.container_for_logs_sliding_window.Len())
		//fmt.Println("actual length", a.container_for_logs_sliding_window.Len())
	}
	return a.container_for_logs_sliding_window.Len() > a.threshold_requests_per_second*a.time_interval_for_average
}

// TODO when to use * and when to use a copy ?
func (a *Alarm) GetAlarmState() bool {
	return a.total_hits > a.threshold_requests_per_second*a.time_interval_for_average
}

func (a *Alarm) GetAverageRequestPerSecond() float64 {
	return float64(a.total_hits) / float64(a.time_interval_for_average)
}

func (a *Alarm) GenerateAlarmMsg(timestamp int) string {
	triggeredTime := time.Unix(int64(timestamp), 0).UTC()
	alertMessage := fmt.Sprintf("High traffic generated an alert - hits = %.2f, triggered at %s", a.GetAverageRequestPerSecond(), triggeredTime.Format(time.RFC3339))
	return alertMessage
}

func (a *Alarm) GenerateRecoveryAlarmMsg(timestamp int) string {
	triggeredTime := time.Unix(int64(timestamp), 0).UTC()
	alertMessage := fmt.Sprintf("Traffic has recovered - hits = %.2f, recovered at %s", a.GetAverageRequestPerSecond(), triggeredTime.Format(time.RFC3339))
	return alertMessage
}

func (a *Alarm) resetOldEntries(timestamp int) {
	for e := a.container_for_logs_sliding_window.Front(); e != nil; {
		next := e.Next()
		e_value := e.Value.(AlarmEntry)

		if e_value.time_in_seconds+a.time_interval_for_average < timestamp {
			a.total_hits = a.total_hits - e_value.qty
			a.container_for_logs_sliding_window.Remove(e)
		} else {
			break
		}
		e = next
	}
}
