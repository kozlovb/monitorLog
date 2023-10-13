package main

import (
	"monitorLog/parser"
	"monitorLog/stats"
	"time"
)

// todo alarm to alert :) rename
type HttpLogMonitor struct {
	time_interval_stats           time.Duration
	time_interval_traffic_average time.Duration
	threshold_traffic_alert       int
	alert                         *Alarm
	stats                         *stats.Statistics
	report_display_chan           chan *stats.Report
	alert_display_chan            chan *string
	display                       Display
}

// TODO when *HttpLogMonitor and when not
func NewHttpLogMonitor(time_interval_stats time.Duration, time_interval_traffic_average time.Duration, threshold_traffic_alert int) *HttpLogMonitor {
	//TODO remove this channels from monitor can be in a display contructor
	report_display_chan := make(chan *stats.Report)
	alert_chan := make(chan *string)

	return &HttpLogMonitor{
		time_interval_stats:           time_interval_stats,
		time_interval_traffic_average: time_interval_traffic_average,
		threshold_traffic_alert:       threshold_traffic_alert,
		alert:                         NewAlarm(int(time_interval_traffic_average.Seconds()), threshold_traffic_alert),
		stats:                         stats.NewStatistics(),
		report_display_chan:           report_display_chan,
		alert_display_chan:            alert_chan,
		display:                       Display{report_chan: report_display_chan, alert_chan: alert_chan},
	}
}

func (h *HttpLogMonitor) Start(log_file_name *string) {
	return_channel := make(chan *parser.Entity)
	r := Reader{file_name: log_file_name, return_channel: return_channel, parser: parser.NewParser()}
	go r.Read()
	go h.run(return_channel)
	h.display.debug_display()
}

func (h *HttpLogMonitor) Stop() {
	// TODO may be like a stop channel
}

// todo specify that its input only
func (h *HttpLogMonitor) run(read_channel chan *parser.Entity) {
	// can be customized to return average for less than 2 seconds

	alarm_state := false
	done := make(chan bool, 1)
	//relative_log_file_time := 0 //timein log can run faster than in a file
	iterationIndex := 0
	previous_report_time := 0
	line_number := 0
	//TODO start main loop after creating channels but before providign them to classes
	previous_alert_state := false
	for {
		select {
		case c := <-read_channel:
			//fmt.Println("lne number ", line_number)
			line_number += 1
			//alarm state class field ?

			h.alert.RegisterEntry(c.Timestamp)
			alarm_state = h.alert.GetAlarmState()
			if previous_report_time == 0 {
				previous_report_time = c.Timestamp
			}
			//alarm to alert
			if alarm_state != previous_alert_state {

				if alarm_state {
					a := h.alert.generateAlarmMsg(c.Timestamp)
					h.display.alert_chan <- &a
				} else {
					a := h.alert.generateRecoveryAlarmMsg(c.Timestamp)
					h.display.alert_chan <- &a
				}

				previous_alert_state = alarm_state
			}
			if c.Timestamp > previous_report_time+int(time_interval_stats.Seconds()) {
				h.display.report_chan <- h.stats.Report()
				h.stats.Clear()
				previous_report_time = c.Timestamp
			}

			h.stats.RegisterEntry(c)
			iterationIndex++
			if iterationIndex >= 4828 {
				done <- true
			}
		case <-done:
			return
		}
	}
}
