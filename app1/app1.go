package app1

import (
	"fmt"
	"net/http"

	"github.com/linuxexam/goweb/router"
)

func init() {
	router.RegisterApp("app1", http.HandlerFunc(app1Handler))
}

func app1Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "I am app1.")
}
