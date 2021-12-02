package cert

///////////////////////////////////////////////////////////////////////////////////////
// This file provides the web app and registers itself via init().

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/linuxexam/goweb/router"
)

//embed resource files
//go:embed *.html *.png
var webUI embed.FS

var (
	appName = "CertChecker"
	tpls    = template.Must(template.ParseFS(webUI, "*.html"))
)

type PageCert struct {
	Title       string
	UrlsToCheck string
	Result      string
}

func certHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/"+appName+"/static/") {
		fs := http.StripPrefix("/"+appName+"/static/", http.FileServer(http.FS(webUI)))
		fs.ServeHTTP(w, r)
		return
	}
	_certHandler(w, r)
}
func _certHandler(w http.ResponseWriter, r *http.Request) {
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
			r, _ := CheckCert(url)
			r = "==========" + url + "==========\n" + r
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

func init() {
	router.RegisterApp(appName, http.HandlerFunc(certHandler))
}
