package main

import (
	"log"

	_ "github.com/linuxexam/goweb/cert"
	"github.com/linuxexam/goweb/router"
)

func main() {
	log.Fatal(router.Run(":8080"))
}
