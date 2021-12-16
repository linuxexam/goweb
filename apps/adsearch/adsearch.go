// Package adsearch implements simple functions to search from AD as ldap.
package adsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

type Session struct {
	LdapURL string `json:"ldapURL"`
	BindDN  string `json:"bindDN"`
	BindPWD string `json:"bindPwd"`
	BaseDN  string `json:"baseDN"`
	conn    *ldap.Conn
}

func NewSession(ldapURL, bindDN, bindPWD, baseDN string) (*Session, error) {
	conn, err := ldap.DialURL(ldapURL)
	if err != nil {
		return nil, err
	}

	err = conn.Bind(bindDN, bindPWD)
	if err != nil {
		return nil, err
	}
	return &Session{
		LdapURL: ldapURL,
		BindDN:  bindDN,
		BindPWD: bindPWD,
		BaseDN:  baseDN,
		conn:    conn,
	}, nil
}

func (s *Session) Close() {
	s.conn.Close()
}

func (s *Session) FindEmail(DN string) (string, error) {
	CN := strings.Split(strings.SplitN(DN, ",", 2)[0], "=")[1]
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(CN))
	attrs := []string{"mail"}
	searchReq := ldap.NewSearchRequest(s.BaseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false,
		filter, attrs, []ldap.Control{})

	r, err := s.conn.Search(searchReq)
	if err != nil {
		return "", err
	}

	if len(r.Entries) == 0 {
		return "", errors.New("not found")
	}

	entry := r.Entries[0]
	if len(entry.Attributes) == 0 {
		return "", errors.New("no email")
	}
	return entry.Attributes[0].Values[0], nil
}
func (s *Session) FindSAM(DN string) (string, error) {
	CN := strings.Split(strings.SplitN(DN, ",", 2)[0], "=")[1]
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(CN))
	attrs := []string{"sAMAccountName"}
	searchReq := ldap.NewSearchRequest(s.BaseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false,
		filter, attrs, []ldap.Control{})

	r, err := s.conn.Search(searchReq)
	if err != nil {
		return "", err
	}

	if len(r.Entries) == 0 {
		return "", errors.New("not found")
	}

	entry := r.Entries[0]
	if len(entry.Attributes) == 0 {
		return "", errors.New("no info")
	}
	return entry.Attributes[0].Values[0], nil
}
func (s *Session) FindManager(DN string) (string, error) {
	CN := strings.Split(strings.SplitN(DN, ",", 2)[0], "=")[1]
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(CN))
	attrs := []string{"manager"}
	searchReq := ldap.NewSearchRequest(s.BaseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false,
		filter, attrs, []ldap.Control{})

	r, err := s.conn.Search(searchReq)
	if err != nil {
		return "", err
	}

	if len(r.Entries) == 0 {
		return "", errors.New("not found")
	}

	entry := r.Entries[0]
	if len(entry.Attributes) == 0 {
		return "", errors.New("no manager")
	}
	return entry.Attributes[0].Values[0], nil
}

func (s *Session) FindManagers(DN string) []string {
	var mgrs []string

	mgr := DN
	var err error
	for {
		mgr, err = s.FindManager(mgr)
		if err != nil {
			break
		}
		mgrs = append(mgrs, mgr)
	}
	return mgrs
}

func (s *Session) FindGroups(DN string) (grps []string, err error) {
	CN := strings.Split(strings.SplitN(DN, ",", 2)[0], "=")[1]
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(CN))
	attrs := []string{"memberOf"}
	searchReq := ldap.NewSearchRequest(s.BaseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false,
		filter, attrs, []ldap.Control{})
	r, err := s.conn.Search(searchReq)
	if err != nil {
		return nil, err
	}
	if len(r.Entries) == 0 {
		return nil, errors.New("not found")
	}
	if len(r.Entries[0].Attributes) == 0 {
		return nil, errors.New("no group found")
	}
	return r.Entries[0].Attributes[0].Values, nil
}

func (s *Session) FindDNBySAM(sam string) (string, error) {
	filter := fmt.Sprintf("(sAMAccountName=%s)", ldap.EscapeFilter(sam))
	attrs := []string{"DN"}
	searchReq := ldap.NewSearchRequest(s.BaseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false,
		filter, attrs, []ldap.Control{})
	r, err := s.conn.Search(searchReq)
	if err != nil {
		return "", err
	}
	if len(r.Entries) == 0 {
		return "", errors.New("not found")
	}
	return r.Entries[0].DN, nil
}

func (s *Session) FindDNsByMail(mail string) ([]string, error) {
	filter := fmt.Sprintf("(mail=%s)", ldap.EscapeFilter(mail))
	attrs := []string{"DN"}
	searchReq := ldap.NewSearchRequest(s.BaseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false,
		filter, attrs, []ldap.Control{})
	r, err := s.conn.Search(searchReq)
	if err != nil {
		return nil, err
	}
	if len(r.Entries) == 0 {
		return nil, errors.New("not found")
	}
	var DNs []string
	for _, entry := range r.Entries {
		DNs = append(DNs, entry.DN)
	}
	return DNs, nil
}
func (s *Session) FindDNsByName(givenName, surName string) ([]string, error) {
	filter := fmt.Sprintf("(&(givenName=%s)(sn=%s))", ldap.EscapeFilter(givenName), ldap.EscapeFilter(surName))
	attrs := []string{"DN"}
	searchReq := ldap.NewSearchRequest(s.BaseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false,
		filter, attrs, []ldap.Control{})
	r, err := s.conn.Search(searchReq)
	if err != nil {
		return nil, err
	}
	if len(r.Entries) == 0 {
		return nil, errors.New("not found")
	}
	var DNs []string
	for _, entry := range r.Entries {
		DNs = append(DNs, entry.DN)
	}
	return DNs, nil
}

func NewSessionFromJson(jsonFile string) (*Session, error) {
	jsonBytes, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	var s Session
	err = json.Unmarshal(jsonBytes, &s)
	if err != nil {
		return nil, err
	}

	conn, err := ldap.DialURL(s.LdapURL)
	if err != nil {
		return nil, err
	}

	err = conn.Bind(s.BindDN, s.BindPWD)
	if err != nil {
		return nil, err
	}

	s.conn = conn
	return &s, nil
}
