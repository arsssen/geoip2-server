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
		WorkDir      string `arg:"env:WORK_DIR, -W, --work-dir" default:""`
	}
	arg.MustParse(&args)

	resolver, err := app.NewResolver(context.Background(), app.Opts{
		LicenseKey:   args.LicenseKey,
		DBFile:       args.DBFile,
		UpdatePeriod: args.UpdatePeriod,
		Logger:       log.Printf,
		WorkingDir:   args.WorkDir,
	})

	// example getting ip directly, without http handler:
	// note that GetLocationJSON method will work only after db file is loaded, if db is not ready yet, it will return error
	//ipStr := net.ParseIP("8.8.8.8")
	//locShort, _ := resolver.GetLocationJSON(ipStr, "short")
	//locFull, _ := resolver.GetLocationJSON(ipStr, "")
	//log.Printf("short: %s, full: %s", locShort, locFull)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(resolver.StartHTTP(args.ListenAddr))
}
