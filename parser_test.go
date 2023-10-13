package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func Test_parseLogString(t *testing.T) {

	parser := NewParser()

	type args struct {
		input string
	}
	tests := []struct {
		name           string
		input          string
		expected       *Entity
		expected_error error
	}{
		{"valid_log_entry",
			"\"10.0.0.2\",\"-\",\"apache\",1549573859,\"GET /api/user HTTP/1.0\",200,1234",
			&Entity{
				ip_address:       "10.0.0.2",
				user_identifier:  "-",
				remote_user_name: "apache",
				timestamp:        1549573859,
				request:          "GET /api/user HTTP/1.0",
				section:          "/api",
				http_status_code: 200,
				response_size:    1234,
			},
			nil,
		},
		{"missing_user_identifier",
			"\"10.0.0.2\",\"apache\",1549573859,\"GET /api/user HTTP/1.0\",200,1234",
			nil,
			errors.New("access log line '\"10.0.0.2\",\"apache\",1549573859,\"GET /api/user HTTP/1.0\",200,1234' does not match given format '^\"(?P<remote_addr>[^\"]*)\",\"(?P<user_identifier>[^\"]*)\",\"(?P<remote_user>[^\"]*)\",(?P<time_local>[^,]*),\"(?P<request>[^\"]*)\",(?P<status>[^,]*),(?P<bytes_sent>[^ ]*)'"),
		},
		{"wrong_request_format",
			"\"10.0.0.2\",\"-\",\"apache\",1549573859,\"GET/api/user HTTP/1.0\",200,1234",
			nil,
			errors.New("Invalid structure of request part: [GET/api/user HTTP/1.0] of log string: \"10.0.0.2\",\"-\",\"apache\",1549573859,\"GET/api/user HTTP/1.0\",200,1234"),
		},
		{"wrong_section_format",
			"\"10.0.0.2\",\"-\",\"apache\",1549573859,\"GET api/user HTTP/1.0\",200,1234",
			nil,
			errors.New("Error parsing a request field for a log string \"10.0.0.2\",\"-\",\"apache\",1549573859,\"GET api/user HTTP/1.0\",200,1234, error: bad format of resource api/user"),
		},
		{"timestamp_is_not_an_int",
			"\"10.0.0.2\",\"-\",\"apache\",15AB49573859,\"GET api/user HTTP/1.0\",200,1234",
			nil,
			errors.New("Error parsing a request field for a log string \"10.0.0.2\",\"-\",\"apache\",15AB49573859,\"GET api/user HTTP/1.0\",200,1234, error: bad format of resource api/user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err_actual := parser.ParseLogString(tt.input)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("parseLogString \nactual: %v \nexpected: %v", actual, tt.expected)
			}
			if err_actual == nil && tt.expected_error != nil {
				t.Errorf("parseLogString \nactual_error: nil \nexpected_error: %v", tt.expected)
			}
			if err_actual != nil && tt.expected_error == nil {
				t.Errorf("parseLogString \nactual_error: %v \nexpected_error: nil", err_actual)
			}
			if (err_actual != nil && tt.expected_error != nil) && !strings.Contains(err_actual.Error(), tt.expected_error.Error()) {
				t.Errorf("parseLogString \nactual_error: %v \nexpected_error: %v", err_actual, tt.expected_error)
			}
		})
	}
}

//checks for negative numbers
// checks for realistic ip etc ...
// code convention remark
