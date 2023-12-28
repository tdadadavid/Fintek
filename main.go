package main

import (
	"github/tdadadavid/fingreat/api"
)

func main() {
	server := api.NewServer(".")
	server.Start(8080)
}