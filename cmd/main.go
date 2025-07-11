package main

import (
	"log/slog"
	"owe-flow/internal/gapi"
)

func main() {
	slog.Info("Running Owe-Flow")
	gapi.ReadSheed()
}
