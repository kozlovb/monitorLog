package parser

import (
	"fmt"

	"strconv"
	"strings"

	"github.com/satyrius/gonx"
)

const (
	format = `"$remote_addr","$user_identifier","$remote_user",$time_local,"$request",$status,$bytes_sent`
)

type Parser struct {
	gonx_parser *gonx.Parser
}

type Entity struct {
	Ip_address       string
	User_identifier  string
	Remote_user_name string
	Timestamp        int
	Request          string
	Section          string
	Http_status_code int
	Response_size    int
}

func NewParser() *Parser {
	return &Parser{
		gonx_parser: gonx.NewParser(format),
	}
}

// This func parses the following log entiries
// 10.0.0.4","-","apache",1549574334,"POST /report HTTP/1.0",404,1307
func (p *Parser) ParseLogString(log_string string) (*Entity, error) {

	entry, err := p.gonx_parser.ParseString(log_string)
	if err != nil {
		return nil, err
	}

	remote_address, err := entry.Field("remote_addr")
	if err != nil {
		fmt.Printf("Error parcing a remote_addr field for a log string %s, error: %v\n", log_string, err)
		return nil, err
	}

	user_identifier, err := entry.Field("user_identifier")
	if err != nil {
		fmt.Printf("Error parsing a user_identifier field for a log string %s, error: %v\n", log_string, err)
		return nil, err
	}

	remote_user_name, err := entry.Field("remote_user")
	if err != nil {
		fmt.Printf("Error parsing a remote_user field for a log string %s, error: %v\n", log_string, err)
		return nil, err
	}

	timestamp, err := entry.Field("time_local")
	if err != nil {
		fmt.Printf("Error parsing a time_local field for a log string %s, error: %v\n", log_string, err)
		return nil, err
	}

	request, err := entry.Field("request")
	if err != nil {
		fmt.Printf("Error parsing a request field for a log string %s, error: %v\n", log_string, err)
		return nil, err
	}

	request_detailed := strings.Fields(request)

	if len(request_detailed) != 3 {
		return nil, fmt.Errorf("Invalid structure of request part: %s of log string: %s", request_detailed, log_string)
	}

	path := request_detailed[1]
	section, err := getSectionFromPath(path)
	if err != nil {
		return nil, fmt.Errorf("Error parsing a request field for a log string %s, error: %v\n", log_string, err)
	}

	statusCode, err := entry.Field("status")
	if err != nil {
		fmt.Printf("Error parsing a status field for a log string %s, error: %v\n", log_string, err)
		return nil, err
	}

	responseSize, err := entry.Field("bytes_sent")
	if err != nil {
		fmt.Printf("Error parsing a bytes_sent field for a log string %s, error: %v\n", log_string, err)
		return nil, err
	}

	timestamp_int, err := strconv.Atoi(timestamp)
	if err != nil {
		fmt.Printf("Couldn't convert to int a timestamp %s of a log string %s, error: %v\n", timestamp, log_string, err)
		return nil, err
	}

	http_status_code_int, err := strconv.Atoi(statusCode)
	if err != nil {
		fmt.Printf("Couldn't convert to int a statusCode %s of a log string %s, error: %v\n", statusCode, log_string, err)
		return nil, err
	}

	response_size_int, err := strconv.Atoi(responseSize)
	if err != nil {
		fmt.Printf("Couldn't convert to int a responseSize %s of a log string %s, error: %v\n", responseSize, log_string, err)
		return nil, err
	}

	return &Entity{Ip_address: remote_address,
		User_identifier:  user_identifier,
		Remote_user_name: remote_user_name,
		Timestamp:        timestamp_int,
		Request:          request,
		Section:          section,
		Http_status_code: http_status_code_int,
		Response_size:    response_size_int}, nil
}

// getSectionFromPath returns a section from a resource path. If resource doesn't start
// with "/" the error is returned if there is more than one "/" then the content before the
// second "/" is returned. If there is only one "/"then he whole string is a section.
// "/api/report" returns "/api"
// "/api" returns "/api"
func getSectionFromPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '/' {
		return "", fmt.Errorf("bad format of resource %s", path)
	}
	split_groups := strings.Split(path, "/")
	if len(split_groups) == 1 {
		return "/" + split_groups[0], nil
	}
	return "/" + split_groups[1], nil
}
