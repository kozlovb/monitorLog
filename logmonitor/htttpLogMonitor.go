package logmonitor

import (
	"monitorLog/alert"
	"monitorLog/display"
	"monitorLog/parser"
	"monitorLog/reader"
	"monitorLog/stats"
	"time"
)

// todo alert to alert :) rename
type HttpLogMonitor struct {
	time_interval_stats           time.Duration
	time_interval_traffic_average time.Duration
	threshold_traffic_alert       int
	alert                         *alert.Alert
	stats                         *stats.Statistics
	report_display_chan           chan *stats.Report
	alert_display_chan            chan *string
	display                       display.Display
}

func NewHttpLogMonitor(time_interval_stats time.Duration, time_interval_traffic_average time.Duration, threshold_traffic_alert int) *HttpLogMonitor {
	//TODO remove this channels from monitor can be in a display contructor
	report_display_chan := make(chan *stats.Report)
	alert_chan := make(chan *string)

	return &HttpLogMonitor{
		time_interval_stats:           time_interval_stats,
		time_interval_traffic_average: time_interval_traffic_average,
		threshold_traffic_alert:       threshold_traffic_alert,
		alert:                         alert.NewAlert(int(time_interval_traffic_average.Seconds()), threshold_traffic_alert),
		stats:                         stats.NewStatistics(),
		report_display_chan:           report_display_chan,
		alert_display_chan:            alert_chan,
		display:                       display.Display{Report_chan: report_display_chan, Alert_chan: alert_chan},
	}
}

func (h *HttpLogMonitor) Start(log_file_name *string) {
	return_channel := make(chan *parser.Entity)
	r := reader.Reader{File_name: log_file_name, Return_channel: return_channel, Parser: parser.NewParser()}
	go r.Read()
	go h.run(return_channel)
	h.display.Display()
}

// todo specify that its input only
func (h *HttpLogMonitor) run(read_channel chan *parser.Entity) {
	// can be customized to return average for less than 2 seconds
	alert_state := false
	//relative_log_file_time := 0 //timein log can run faster than in a file
	previous_report_time := 0

	//TODO start main loop after creating channels but before providign them to classes
	previous_alert_state := false
	for {
		select {
		case c, ok := <-read_channel:

			if !ok {
				return
			}

			h.alert.RegisterEntry(c.Timestamp)
			alert_state = h.alert.GetAlertState()
			if previous_report_time == 0 {
				previous_report_time = c.Timestamp
			}
			if alert_state != previous_alert_state {

				if alert_state {
					a := h.alert.GenerateAlertMsg(c.Timestamp)
					h.display.Alert_chan <- &a
				} else {
					a := h.alert.GenerateRecoveryAlertMsg(c.Timestamp)
					h.display.Alert_chan <- &a
				}
				previous_alert_state = alert_state
			}
			if c.Timestamp > previous_report_time+int(h.time_interval_stats.Seconds()) {
				h.display.Report_chan <- h.stats.Report()
				h.stats.Clear()
				previous_report_time = c.Timestamp
			}
			h.stats.RegisterEntry(c)
		}
	}
}
