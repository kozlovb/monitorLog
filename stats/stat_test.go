package stats

import (
	"monitorLog/parser"
	"reflect"
	"testing"
)

func Test_retriveReport(t *testing.T) {

	stats := NewStatistics()

	// entries are not sorted by time
	entities := [...]parser.Entity{{
		Ip_address:       "10.0.0.2",
		User_identifier:  "-",
		Remote_user_name: "apache",
		Timestamp:        1549573859,
		Request:          "GET /api/user HTTP/1.0",
		Section:          "/api",
		Http_status_code: 200,
		Response_size:    1234,
	},
		{
			Ip_address:       "10.0.0.1",
			User_identifier:  "-",
			Remote_user_name: "apache",
			Timestamp:        1549573860,
			Request:          "GET /api/user HTTP/1.0",
			Section:          "/report",
			Http_status_code: 200,
			Response_size:    1234,
		},
		{
			Ip_address:       "10.0.0.2",
			User_identifier:  "-",
			Remote_user_name: "apache",
			Timestamp:        1549573859,
			Request:          "GET /api/user HTTP/1.0",
			Section:          "/api",
			Http_status_code: 200,
			Response_size:    1234,
		}}

	for _, entity := range entities {
		stats.RegisterEntry(&entity)
	}

	report_actual := stats.Report()

	report_expected := &Report{
		Number_of_hits:   2,
		Section:          "/api",
		Ip_from:          "10.0.0.2",
		Hits_from_max_ip: 2,
	}

	if !reflect.DeepEqual(*report_actual, *report_expected) {
		t.Errorf("Report \nactual: %v \nexpected: %v", report_actual, report_expected)
	}

}
