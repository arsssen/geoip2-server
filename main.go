package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/oschwald/maxminddb-golang"
)

const downloadURL = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz"
const shaURL = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz.sha256"
const downloadPath = "./downloaded.tgz"
const dbFileName = "GeoLite2-City.mmdb" // file name inside the downloaded tgz archive

var args struct {
	LicenseKey   string `arg:"env:LICENSE_KEY, -L, --license-key" `
	ListenAddr   string `arg:"env:LISTEN_ADDR, -A, --listen-addr" default:":8080"`
	DBFile       string `arg:"env:DB_FILE, -D, --db-file"`
	UpdatePeriod string `arg:"env:UPDATE_PERIOD, -U, --update-period" default:"240h"`
}

func main() {

	arg.MustParse(&args)
	s := newServer(args.ListenAddr)
	setDB := func(filename string) {
		log.Printf("opening db file %s", filename)
		db, err := maxminddb.Open(filename)
		if err != nil {
			log.Printf("cannot open db %s: %s", filename, err)
			return
		}
		s.setDB(db)
	}
	if args.DBFile != "" {
		setDB(args.DBFile)
	}
	if args.LicenseKey != "" {
		period, e := time.ParseDuration(args.UpdatePeriod)
		if e != nil {
			log.Fatalf("invalid update period: %s", args.UpdatePeriod)
		}
		go startUpdater(period, setDB)
	}
	if args.DBFile == "" && args.LicenseKey == "" {
		log.Fatal("cannot run without license key or existing db file")
	}

	s.start()

}

func startUpdater(period time.Duration, update func(string)) {
	log.Printf("starting updater (%s)", period)
	t := time.NewTicker(period)
	lastSChecksum := ""

	for ; true; <-t.C {
		sha, e := getContent(fmt.Sprintf(shaURL, args.LicenseKey))
		if e != nil {
			log.Printf("cannot get sha: %s", e)
			continue
		}
		log.Printf("sha: %s", sha)
		if sha == lastSChecksum {
			log.Print("no update")
			continue
		}
		lastSChecksum = sha
		if e = downloadFile(downloadPath, fmt.Sprintf(downloadURL, args.LicenseKey)); e != nil {
			log.Printf("download error: %s", e)
			continue
		}
		newFile, e := extractTarGz(downloadPath, dbFileName)
		if e != nil {
			log.Printf("extract error: %s", e)
			continue
		}
		update(newFile)
	}
}
