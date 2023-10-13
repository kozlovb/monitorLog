package main

// reader parser ?
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/satyrius/gonx"
)

const (
	// put example here
	format2 = `$remote_addr,$user_identifier,$remote_user,$time_local,"$request",$status,$bytes_sent`
)

//"remotehost","rfc931","authuser","date","request","status","bytes"
//"10.0.0.2","-","apache",1549573860,"GET /api/user HTTP/1.0",200,1234

// Who passes channel to who ?
// take parcer from the other solution

// introduce tailermay be later
// extend to take parser as an argument probably an interface
// TODO thouyghts if file is slow to get read out
type Reader struct {
	//parser   *logparser.HTTPd
	file_name      *string
	return_channel chan *Entity
}

// tODO rename to reader ?
type Entity struct {
	ip_address       string
	user_identifier  string
	remote_user_name string
	timestamp        int
	request          string
	section          string
	http_status_code int
	response_size    int
}

func (m *Reader) Read() {
	// Create a parser
	parser := gonx.NewParser(format2)
	file, _ := os.Open(*m.file_name)
	//if err != nil
	//{
	// log.Fatal(err)
	//}
	defer file.Close()
	defer close(m.return_channel)

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example

	for scanner.Scan() {

		//fmt.Println(scanner.Text())
		logLine := scanner.Text()
		//will this tring survive ?

		// Parse the log line
		entry, err := parser.ParseString(logLine)
		if err != nil {
			fmt.Println("Error parsing log line:", err)
			//TODO ignore first lne
			continue
		}

		//TODO - what to do with a second entry ?

		// Access parsed fields
		ipAddress, _ := entry.Field("remote_addr")
		user_identifier, _ := entry.Field("user_identifier")
		remote_user_name, _ := entry.Field("remote_user")
		timestamp, _ := entry.Field("time_local")
		request, _ := entry.Field("request")
		request_detailed := strings.Fields(request)

		//10.0.0.4","-","apache",1549574334,"POST /report HTTP/1.0",404,1307

		//TODO print error
		if len(request_detailed) != 3 {
			continue
			//return fmt.Errorf("Invalid structure of request part: %s", request_detailed)
		}

		//r.Method = parts[0]
		path := request_detailed[1]
		section, _ := getSectionFromResource(path)
		// TODO error tretment
		//r.Protocol = parts[2]

		statusCode, _ := entry.Field("status")
		responseSize, _ := entry.Field("body_bytes_sent")

		// Print parsed fields
		// `$remote_addr,$user_identifier,$remote_user,$time_local,"$request",$status,$bytes_sent`

		timestamp_int, _ := strconv.Atoi(timestamp)
		http_status_code_int, _ := strconv.Atoi(statusCode)
		response_size_int, _ := strconv.Atoi(responseSize)

		m.return_channel <- &Entity{ip_address: ipAddress,
			user_identifier:  user_identifier,
			remote_user_name: remote_user_name,
			timestamp:        timestamp_int,
			request:          request,
			section:          section,
			http_status_code: http_status_code_int,
			response_size:    response_size_int}

	}
	/*
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}*/
}

// getSectionFromResource returns a section from a resource path.
// A section is defined as being what's before the second '/' in the resource section.
// Eg. the section for '/pages/create' is '/pages'.
// Returns an error if the path is the empty string
func getSectionFromResource(path string) (string, error) {

	// TODO rewrite this shit to avoid plagiat copmpain
	if !strings.HasPrefix(path, "/") { // Reject paths that don't start with /
		return "", fmt.Errorf("cannot get section from path %s", path)
	}
	// Remove leading "/" since I'm sure it's there
	stripped := strings.TrimLeft(path, "/")
	// Split on middle "/"s and take the first
	split := strings.Split(stripped, "/")
	return "/" + split[0], nil
}
