package cert

import (
	"fmt"
	"testing"
)

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
			r, err := PrintCert(url, "json")
			if err != nil {
				fmt.Println(err)
			}

		}(url)
	}

	for range urls {
		fmt.Print(<-chOut)
	}
}
