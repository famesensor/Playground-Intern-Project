package main

import (
	"log"

	"github.com/HangoKub/Hango-service/protocol"
)

func main() {
	if err := protocol.ServeHttp(); err != nil {
		log.Fatal(err)
	}
}
