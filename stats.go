package main

import (
	"fmt"
	"monitorLog/parser"
)

type SectionMapEntry struct {
	number_of_hits  int
	indexes_entries []int
}

type Report struct {
	number_of_hits   int
	section          string
	ip_from          string
	hits_from_max_ip int
}

type Statistics struct {
	hits_map           map[string]*SectionMapEntry
	curent_max_section string
	current_max_hits   int
	entries            []*parser.Entity
}

func NewStatistics() *Statistics {
	return &Statistics{
		hits_map:           make(map[string]*SectionMapEntry),
		curent_max_section: "",
		current_max_hits:   0,
		entries:            make([]*parser.Entity, 0),
	}
}

// what this is supposed to do ?

func (s *Statistics) RegisterEntry(entry *parser.Entity) {
	//how to register en entry ?
	s.entries = append(s.entries, entry)
	last_entry_index := len(s.entries) - 1

	new_count := 0
	if sectionEntry, ok := s.hits_map[entry.Section]; ok {
		new_count = sectionEntry.number_of_hits + 1
		sectionEntry.number_of_hits = new_count

		sectionEntry.indexes_entries = append(sectionEntry.indexes_entries, last_entry_index)
	} else {
		new_count = 1
		s.hits_map[entry.Section] = &SectionMapEntry{number_of_hits: new_count, indexes_entries: []int{last_entry_index}}
	}
	if new_count > s.current_max_hits {
		s.curent_max_section = entry.Section
		s.current_max_hits = new_count
	}
}

func (s *Statistics) Report() *Report {
	ip_from_max_hits, max_hits_from_this_ip := s.findMaxIP(&s.curent_max_section)
	return &Report{number_of_hits: s.current_max_hits, section: s.curent_max_section, ip_from: ip_from_max_hits, hits_from_max_ip: max_hits_from_this_ip}
}

// TODO when * and when copy of a class ?
func (s *Statistics) findMaxIP(section *string) (string, int) {
	section_entry, ok := s.hits_map[*section]
	if !ok {
		fmt.Print("trying to get the ... ")
		return "None", -1
	}
	current_max_ip := ""
	current_max_hits := 0
	ip_map := make(map[string]int)
	//for index, value := range numbers {
	for _, index_in_entries := range section_entry.indexes_entries {
		ip := s.entries[index_in_entries].Ip_address
		new_req_from := 0
		if requests_from, ok := ip_map[ip]; ok {
			//TODO check in debugger if it updates the map
			new_req_from = requests_from + 1
		} else {
			new_req_from = 1
		}
		ip_map[ip] = new_req_from
		if new_req_from > current_max_hits {
			current_max_hits = new_req_from
			current_max_ip = ip
		}

	}
	return current_max_ip, current_max_hits
}

func (s *Statistics) Clear() {
	s.hits_map = make(map[string]*SectionMapEntry)
	s.curent_max_section = "invalid" // todo  a const
	s.current_max_hits = 0
}
