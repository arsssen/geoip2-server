package app

import (
	"context"
	"fmt"
	"time"
)

const downloadURL = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz"
const shaURL = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz.sha256"
const downloadPath = "./downloaded.tgz"
const dbFileName = "GeoLite2-City.mmdb" // file name inside the downloaded tgz archive

func startUpdater(ctx context.Context, period time.Duration, licenseKey string, update func(string)) {
	log("starting updater (%s)", period)
	t := time.NewTicker(period)
	lastSChecksum := ""

	for ; true; <-t.C {
		select {
		case <-ctx.Done():
			log("stopping updater (ctx cancelled)")
		default:

		}
		sha, e := getContent(fmt.Sprintf(shaURL, licenseKey))
		if e != nil {
			log("cannot get sha: %s", e)
			continue
		}
		log("sha: %s", sha)
		if sha == lastSChecksum {
			log("no update")
			continue
		}
		lastSChecksum = sha
		if e = downloadFile(downloadPath, fmt.Sprintf(downloadURL, licenseKey)); e != nil {
			log("download error: %s", e)
			continue
		}
		newFile, e := extractTarGz(downloadPath, dbFileName)
		if e != nil {
			log("extract error: %s", e)
			continue
		}
		update(newFile)
	}
}
