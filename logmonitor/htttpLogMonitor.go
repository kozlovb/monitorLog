package logmonitor

import (
	"monitorLog/alert"
	"monitorLog/display"
	"monitorLog/parser"
	"monitorLog/reader"
	"monitorLog/stats"
	"time"
)

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

// Creates new HttpLogMonitor struct.
func NewHttpLogMonitor(time_interval_stats time.Duration, time_interval_traffic_average time.Duration, threshold_traffic_alert int) *HttpLogMonitor {
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

// Starts reading from a given log file
// Reading is stopped then the end of the file is reached
// display is stopped by Q key entry
func (h *HttpLogMonitor) Start(log_file_name *string) {
	return_channel := make(chan *parser.Entity)
	r := reader.Reader{File_name: log_file_name, Return_channel: return_channel, Parser: parser.NewParser()}
	go h.run(return_channel)
	go r.Read()
	h.display.Display()
}

// runs the monitor, dispatches data ffrom the read channel to alarm
// and statistics, decides on when to request status from these structs
func (h *HttpLogMonitor) run(read_channel <-chan *parser.Entity) {

	alert_state := false
	previous_report_timestamp := 0

	previous_alert_state := false
	for {
		select {
		case c, ok := <-read_channel:
			if !ok {
				return
			}
			h.alert.RegisterEntry(c.Timestamp)
			alert_state = h.alert.GetAlertState()
			if previous_report_timestamp == 0 {
				previous_report_timestamp = c.Timestamp
			}
			if alert_state != previous_alert_state {
				if alert_state {
					alert_msg := h.alert.GenerateAlertMsg(c.Timestamp)
					h.display.Alert_chan <- &alert_msg
				} else {
					recovery_msg := h.alert.GenerateRecoveryAlertMsg(c.Timestamp)
					h.display.Alert_chan <- &recovery_msg
				}
				previous_alert_state = alert_state
			}

			h.doReports(&previous_report_timestamp, c.Timestamp)
			h.stats.RegisterEntry(c)
		}
	}
}

// Triggers reports generation. Number of report correspond to the time past devided
// by the specified interval
func (h *HttpLogMonitor) doReports(previous_report_timestamp *int, new_timestamp int) {
	interval := int(h.time_interval_stats.Seconds())
	if new_timestamp <= *previous_report_timestamp || interval == 0 {
		return
	}
	ratio := (new_timestamp - *previous_report_timestamp) / int(h.time_interval_stats.Seconds())
	for i := 0; i < ratio; i++ {
		h.display.Report_chan <- h.stats.Report()
		h.stats.Clear()
	}
	if ratio > 0 {
		*previous_report_timestamp = new_timestamp
	}
}
