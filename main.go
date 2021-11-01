package main

import (
	"context"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/arsssen/geoip2-server/app"
)

func main() {
	var args struct {
		LicenseKey   string `arg:"env:LICENSE_KEY, -L, --license-key" `
		ListenAddr   string `arg:"env:LISTEN_ADDR, -A, --listen-addr" default:":8080"`
		DBFile       string `arg:"env:DB_FILE, -D, --db-file"`
		UpdatePeriod string `arg:"env:UPDATE_PERIOD, -U, --update-period" default:"240h"`
	}
	arg.MustParse(&args)
	log.Fatal(app.Run(context.Background(), app.Opts(args), log.Printf))
}
