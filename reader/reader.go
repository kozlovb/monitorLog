package reader

import (
	"bufio"
	"fmt"
	"monitorLog/parser"
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
	File_name      *string
	Return_channel chan *parser.Entity
	Parser         *parser.Parser
}

func (m *Reader) Read() {
	file, err := os.Open(*m.File_name)
	if err != nil {
		fmt.Println("Error opening file a parser", err)
	}
	defer file.Close()
	defer close(m.Return_channel)

	scanner := bufio.NewScanner(file)

	first_line := true
	for scanner.Scan() {
		if first_line {
			first_line = false
			continue
		}
		log_string := scanner.Text()
		parsed_entry, err := m.Parser.ParseLogString(log_string)
		if err != nil {
			fmt.Println(err)
		} else {
			m.Return_channel <- parsed_entry
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error during file reading", err)
	}
}
