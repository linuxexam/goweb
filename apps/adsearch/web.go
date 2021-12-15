package adsearch

import (
	"fmt"
	"net/http"

	"github.com/linuxexam/goweb/router"
)

func init() {
	router.RegisterApp("Active Directory Search", http.HandlerFunc(appHandler))
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "I am app1.")
}
