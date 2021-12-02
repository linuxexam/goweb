package main

import (
	"log"

	_ "github.com/linuxexam/goweb/apps/app1"
	_ "github.com/linuxexam/goweb/apps/cert"
	"github.com/linuxexam/goweb/router"
)

func main() {
	log.Fatal(router.Run(":8080"))
}
