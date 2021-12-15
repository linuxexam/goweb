package adsearch

import (
	"fmt"
	"log"
	"testing"
)

func TestADSearch(t *testing.T) {
	s, err := NewSessionFromJson("bcit-session.json")
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	DNs, err := s.FindDNsByName("jonathan", "zhao")
	if err != nil {
		log.Fatal(err)
	}
	DN := DNs[0]
	mgrs := s.FindManagers(DN)
	for _, mgr := range mgrs {
		mail, err := s.FindEmail(mgr)
		if err != nil {
			fmt.Printf("No email")
			continue
		}
		fmt.Printf("%s\n", mail)
	}

	groups, _ := s.FindGroups(DN)
	for _, grp := range groups {
		fmt.Printf("%s\n", grp)
	}
}
