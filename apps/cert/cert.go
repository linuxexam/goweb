// cert is a lib to parse TLS cert
package cert

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
)

func FindCert(url string) (*x509.Certificate, error) {
	if !strings.Contains(url, ":") {
		url += ":443"
	}
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", url, config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cs := conn.ConnectionState()
	peerCert := cs.PeerCertificates[0]
	return peerCert, nil
}

func JsonCert(cert *x509.Certificate) (string, error) {
	b, err := json.MarshalIndent(cert, " ", " ")
	return string(b), err
}

func YamlCert(cert *x509.Certificate) (string, error) {
	b, err := yaml.Marshal(cert)
	return string(b), err
}

func StringCert(cert *x509.Certificate) string {
	s := new(strings.Builder)
	fmt.Fprintf(s, "subject:    %s\n", cert.Subject.CommonName)
	fmt.Fprintf(s, "SAN:        %s\n", cert.DNSNames)
	fmt.Fprintf(s, "Not before: %s\n", StringTime(cert.NotBefore))
	fmt.Fprintf(s, "Not after:  %s\n", StringTime(cert.NotAfter))
	return s.String()
}

func StringTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func PrintCert(url, format string) (string, error) {
	cert, err := FindCert(url)
	if err != nil {
		return "Not found", err
	}
	switch format {
	case "json":
		return JsonCert(cert)
	case "yaml":
		return YamlCert(cert)
	default:
		return StringCert(cert), nil
	}
}
