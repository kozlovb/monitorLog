package main

import (
	"bufio"
	"fmt"
	"os"
)

//"remotehost","rfc931","authuser","date","request","status","bytes"
//"10.0.0.2","-","apache",1549573860,"GET /api/user HTTP/1.0",200,1234

// Who passes channel to who ?

// introduce tailermay be later
// extend to take parser as an argument probably an interface
// TODO thouyghts if file is slow to get read out
type Reader struct {
	//parser   *logparser.HTTPd
	file_name      *string
	return_channel chan *Entity
	parser         *Parser
}

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

	file, err := os.Open(*m.file_name)
	if err != nil {
		fmt.Println("Error opening file a parser", err)
	}
	defer file.Close()
	defer close(m.return_channel)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		log_string := scanner.Text()
		parsed_entry, err := m.parser.ParseLogString(log_string)
		if err != nil {
			fmt.Println(err)
		} else {
			m.return_channel <- parsed_entry
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error during file reading", err)
	}
}
