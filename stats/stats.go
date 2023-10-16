package stats

import (
	"fmt"
	"monitorLog/parser"
)

type SectionMapEntry struct {
	number_of_hits  int
	indexes_entries []int
}

type Report struct {
	Number_of_hits   int
	Section          string
	Ip_from          string
	Hits_from_max_ip int
}

// In this structure a hash map allows to find SectionMapEntry that correspondss to
// a given section. SectionMapEntry stores all the indexes in entries slice for this section
// so that all supplementary data can be rerieved.
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

// Registers a new entry to the statistics struct.
func (s *Statistics) RegisterEntry(entry *parser.Entity) {
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

// Reports the statistics for the currently registered entries.
func (s *Statistics) Report() *Report {
	ip_from_max_hits, max_hits_from_this_ip := s.findMaxIPbyHits(&s.curent_max_section)
	return &Report{Number_of_hits: s.current_max_hits, Section: s.curent_max_section, Ip_from: ip_from_max_hits, Hits_from_max_ip: max_hits_from_this_ip}
}

// Finds an IP adresse that provided the most hits for the given section.
func (s *Statistics) findMaxIPbyHits(section *string) (string, int) {
	section_entry, ok := s.hits_map[*section]
	if !ok {
		fmt.Printf("findMaxIPbyHits:  section - %s is not in hitmap", *section)
	}
	current_max_ip := ""
	current_max_hits := 0
	ip_map := make(map[string]int)

	for _, index_in_entries := range section_entry.indexes_entries {
		ip := s.entries[index_in_entries].Ip_address
		new_req_from := 0
		if requests_from, ok := ip_map[ip]; ok {
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

// Clears the statistics struct
func (s *Statistics) Clear() {
	s.hits_map = make(map[string]*SectionMapEntry)
	s.curent_max_section = ""
	s.current_max_hits = 0
}
