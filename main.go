package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/linuxexam/goweb/cert"
)

var (
	tpls = template.Must(template.ParseFiles(
		"www/index.html",
		"www/cert.html"))
)

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/cert/", certHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tpls.ExecuteTemplate(w, "index.html", nil)
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////// Certificate Checker /////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
type PageCert struct {
	Title       string
	UrlsToCheck string
	Result      string
}

func certHandler(w http.ResponseWriter, r *http.Request) {
	// test only
	tpls = template.Must(template.ParseFiles(
		"www/index.html",
		"www/cert.html"))

	pd := &PageCert{
		Title:       "Cert Checker",
		UrlsToCheck: "www.bcit.ca",
		Result:      "",
	}
	urlToCheck := r.FormValue("urlToCheck")

	// get and no input
	if urlToCheck == "" && r.Method == "GET" {
		tpls.ExecuteTemplate(w, "cert.html", pd)
		return
	}

	// post
	urls := strings.Fields(urlToCheck)
	sb := new(strings.Builder)
	chOut := make(chan string)
	for _, url := range urls {
		go func(url string) {
			r, _ := cert.CheckCert(url)
			r = "----------" + url + "----------\n" + r
			chOut <- r
		}(url)
	}
	for range urls {
		sb.WriteString(<-chOut)
	}

	pd.UrlsToCheck = urlToCheck
	pd.Result = sb.String()
	tpls.ExecuteTemplate(w, "cert.html", pd)
}
