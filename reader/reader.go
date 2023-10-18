package reader

import (
	"bufio"
	"fmt"
	"monitorLog/parser"
	"os"
)

//"remotehost","rfc931","authuser","date","request","status","bytes"
//"10.0.0.2","-","apache",1549573860,"GET /api/user HTTP/1.0",200,1234

type Reader struct {
	//parser   *logparser.HTTPd
	File_name      *string
	Return_channel chan *parser.Entity
	Parser         *parser.Parser
}

// Reads a file and closes the read_channel at the end
func (r *Reader) Read() {
	file, err := os.Open(*r.File_name)
	if err != nil {
		fmt.Println("Error opening file", err)
	}
	defer file.Close()
	defer close(r.Return_channel)

	scanner := bufio.NewScanner(file)

	first_line := true
	for scanner.Scan() {
		if first_line {
			first_line = false
			continue
		}
		log_string := scanner.Text()
		parsed_entry, err := r.Parser.ParseLogString(log_string)
		if err != nil {
			fmt.Println(err)
		} else {
			r.Return_channel <- parsed_entry
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error during file reading", err)
	}
}
