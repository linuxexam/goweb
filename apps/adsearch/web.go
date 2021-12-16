package adsearch

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/linuxexam/goweb/router"
)

//embed resource files
//go:embed *.html *.png
var webUI embed.FS

var (
	appName = "ActiveDirectorySearch"
	tpls    = template.Must(template.ParseFS(webUI, "*.html"))
)

type PageIndex struct {
	Title    string
	UserName string
	Result   string
}

func init() {
	router.RegisterApp(appName, http.HandlerFunc(appHandler))
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/"+appName+"/static/") {
		fs := http.StripPrefix("/"+appName+"/static/", http.FileServer(http.FS(webUI)))
		fs.ServeHTTP(w, r)
		return
	}
	_appHandler(w, r)
}

func getDNInfo(s *Session, DN string) string {
	var sb strings.Builder
	// basic info
	email, err := s.FindEmail(DN)
	sb.WriteString("Email:\n")

	if err != nil {
		sb.WriteString(err.Error() + "\n")
	}
	sb.WriteString(email + "\n")

	sb.WriteString("\n")
	// managers info
	mgrs := s.FindManagers(DN)

	sb.WriteString("Managers:\n")
	for _, mgr := range mgrs {
		sb.WriteString(mgr + "\n")
	}

	// groups info
	sb.WriteString("\n")
	sb.WriteString("Groups:\n")
	grps, err := s.FindGroups(DN)
	if err != nil {
		sb.WriteString("No groups")
	}
	for _, grp := range grps {
		sb.WriteString(grp + "\n")
	}
	return sb.String()
}

// config file priority order: current work dir, exe's dir
func getConfigFile(filename string) string {

	configFile := filepath.Join("config", "adsearch", filename)

	_, err := os.Stat(configFile)
	if err != nil {
		selfPath, _ := os.Executable()
		selfDir := filepath.Dir(selfPath)
		configFile = filepath.Join(selfDir, "config", "adsearch", filename)
	}

	return configFile
}
func _appHandler(w http.ResponseWriter, r *http.Request) {
	pd := &PageIndex{
		Title:    appName,
		UserName: "Jonathan Zhao",
	}
	username := strings.TrimSpace(r.FormValue("username"))

	// get and no input
	if username == "" && r.Method == "GET" {
		tpls.ExecuteTemplate(w, "index.html", pd)
		return
	}

	// post
	pd.UserName = username

	names := strings.Fields(username)
	if len(names) < 2 {
		pd.Result = "Name format wrong! Correct format is 'GivenName Surname'"
		tpls.ExecuteTemplate(w, "index.html", pd)
		return
	}
	givenName := names[0]
	surname := names[1]

	s, err := NewSessionFromJson(getConfigFile("bcit-session.json"))
	if err != nil {
		pd.Result = err.Error()
		tpls.ExecuteTemplate(w, "index.html", pd)
		return
	}
	defer s.Close()

	DNs, err := s.FindDNsByName(givenName, surname)
	if err != nil {
		pd.Result = err.Error()
		tpls.ExecuteTemplate(w, "index.html", pd)
		return
	}

	n := len(DNs)
	pd.Result += fmt.Sprintf("Found %d records\n", n)

	for i, DN := range DNs {
		pd.Result += fmt.Sprintf("---------------------------------------Record: %d/%d---------------------------------------\n", i+1, n)
		pd.Result += getDNInfo(s, DN)
	}
	tpls.ExecuteTemplate(w, "index.html", pd)
}
