package router

import (
	"embed"
	"html/template"
	"net/http"
	"sync"
)

//go:embed *.html *.png *.css
var webUI embed.FS

var (
	tpls = template.Must(template.ParseFS(webUI, "*.html"))
)
var (
	appsMu sync.RWMutex
	apps   = make(map[string]App)
)

type App struct {
	Name        string
	UrlPattern  string
	HttpHandler http.Handler
}

func RegisterApp(name string, handler http.Handler) {
	appsMu.Lock()
	defer appsMu.Unlock()
	if handler == nil {
		panic("web: nil is not valid handler.")
	}
	if _, dup := apps[name]; dup {
		panic("web: Register called twice for handler " + name)
	}
	app := App{name, "/" + name + "/", handler}
	apps[name] = app
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tpls.ExecuteTemplate(w, "router.html", apps)
}

func Run(addr string) error {
	http.HandleFunc("/", rootHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(webUI))))

	for _, app := range apps {
		http.Handle(app.UrlPattern, app.HttpHandler)
	}
	return http.ListenAndServe(addr, nil)
}
