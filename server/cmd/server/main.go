package main

import (
	"fmt"

	"github.com/Alandres998/go-keeper/server/internal/app/server"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if buildVersion == "" {
		buildVersion = "1.0.0"
	}
	if buildDate == "" {
		buildDate = "2024-12-06"
	}
	if buildCommit == "" {
		buildCommit = "--------"
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	server.RunServer()
}
