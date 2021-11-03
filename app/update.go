package app

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/oschwald/maxminddb-golang"
)

const downloadURL = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz"
const shaURL = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz.sha256"
const downloadPath = "./downloaded.tgz"
const dbFileName = "GeoLite2-City.mmdb" // file name inside the downloaded tgz archive

func (s *resolver) startUpdater(ctx context.Context, period time.Duration, licenseKey string, update func(string)) {
	s.logger("starting updater (%s)", period)
	t := time.NewTicker(period)
	lastSChecksum := ""

	for ; true; <-t.C {
		select {
		case <-ctx.Done():
			s.logger("stopping updater (ctx cancelled)")
		default:

		}
		sha, e := getContent(fmt.Sprintf(shaURL, licenseKey))
		if e != nil {
			s.logger("cannot get sha: %s", e)
			continue
		}
		s.logger("sha: %s", sha)
		if sha == lastSChecksum {
			s.logger("no update")
			continue
		}
		lastSChecksum = sha
		downloadedFile := path.Join(s.workDir, downloadPath)
		if e = s.downloadFile(downloadedFile, fmt.Sprintf(downloadURL, licenseKey)); e != nil {
			s.logger("download error: %s", e)
			continue
		}
		newFile, e := s.extractTarGz(downloadedFile, dbFileName)
		if e != nil {
			s.logger("extract error: %s", e)
			continue
		}
		update(newFile)
	}
}

func (s *resolver) setDB(db *maxminddb.Reader) {
	s.logger("using db: %+v", db.Metadata)
	s.dbMutex.Lock()
	if s.db != nil {
		s.logger("closing previous db: %+v", s.db.Metadata)
		s.db.Close()
	}
	s.db = db
	s.dbMutex.Unlock()
	s.cache.Flush()
}
