package main

import (
	"flag"
	"monitorLog/logmonitor"
	"time"
)

var (
	log_file_name                 = flag.String("log_file_name", "sample_csv.txt", "The path to the log file")
	time_interval_stats           = flag.Duration("time_interval_stats", 10*time.Second, "The positive time interval for metrics calculations")
	time_interval_traffic_average = flag.Duration("time_interval_traffic_average", 120*time.Second, "The positive time interval for traffiic average calculation")
	threshold_traffic_alarm       = flag.Int("threshold_traffic_alarm", 10, "The threshold that triggers a traffic alarm")
)

func main() {
	flag.Parse()
	http_log_monitor := logmonitor.NewHttpLogMonitor(*time_interval_stats, *time_interval_traffic_average, *threshold_traffic_alarm)
	http_log_monitor.Start(log_file_name)
}
