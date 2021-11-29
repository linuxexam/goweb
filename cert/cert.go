package cert

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"
	"testing"
	"time"
)

func CheckCert(url string) (string, error) {
	if !strings.Contains(url, ":") {
		url += ":443"
	}
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", url, config)
	if err != nil {
		return err.Error() + "\n", err
	}
	defer conn.Close()

	cs := conn.ConnectionState()
	peerCert := cs.PeerCertificates[0]
	return sprintX509(peerCert), nil
}

func sprintX509(cert *x509.Certificate) string {
	s := new(strings.Builder)
	fmt.Fprintf(s, "subject:    %s\n", cert.Subject.CommonName)
	fmt.Fprintf(s, "SAN:        %s\n", cert.DNSNames)
	fmt.Fprintf(s, "Not before: %s\n", sprintTime(cert.NotBefore))
	fmt.Fprintf(s, "Not after:  %s\n", sprintTime(cert.NotAfter))
	return s.String()
}

func sprintTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func TestCheckCert(t *testing.T) {
	chOut := make(chan string)
	urls := []string{
		"mail.google.com",
		"www.bcit.ca",
		"notexist.cm",
	}
	for _, url := range urls {
		go func(url string) {
			var r string
			defer func() { chOut <- (url + ":\n" + r) }()
			r, _ = CheckCert(url)
		}(url)
	}

	for range urls {
		fmt.Print(<-chOut)
	}
}
