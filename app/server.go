package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/oschwald/maxminddb-golang"
	"github.com/patrickmn/go-cache"
)

var log func(format string, v ...interface{})

//Opts is the set of the options needed to run the server
type Opts struct {
	LicenseKey   string
	ListenAddr   string
	DBFile       string
	UpdatePeriod string
}

//Run starts the server with given options
func Run(args Opts, loggerFunc ...func(format string, v ...interface{})) error {
	log = func(format string, v ...interface{}) {}
	if len(loggerFunc) > 0 {
		log = loggerFunc[0]
	}
	s := newServer(args.ListenAddr)
	setDB := func(filename string) {
		log("opening db file %s", filename)
		db, err := maxminddb.Open(filename)
		if err != nil {
			log("cannot open db %s: %s", filename, err)
			return
		}
		s.SetDB(db)
	}
	if args.DBFile != "" {
		setDB(args.DBFile)
	}
	if args.LicenseKey != "" {
		period, e := time.ParseDuration(args.UpdatePeriod)
		if e != nil {
			return fmt.Errorf("invalid update period: %s", args.UpdatePeriod)
		}
		go startUpdater(period, args.LicenseKey, setDB)
	}
	if args.DBFile == "" && args.LicenseKey == "" {
		return fmt.Errorf("cannot run without license key or existing db file")
	}

	return s.start()

}

func newServer(listenAddr string) server {
	return server{
		dbMutex:    &sync.RWMutex{},
		listenAddr: listenAddr,
		cache:      cache.New(5*time.Hour, time.Hour),
	}
}

type server struct {
	listenAddr string
	db         *maxminddb.Reader
	dbMutex    *sync.RWMutex
	cache      *cache.Cache
}

func (s *server) start() error {
	http.HandleFunc("/", s.locationHandler)
	log("Starting server at %s", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *server) SetDB(db *maxminddb.Reader) {
	log("using db: %+v", db.Metadata)
	s.dbMutex.Lock()
	if s.db != nil {
		log("closing previous db: %+v", s.db.Metadata)
		s.db.Close()
	}
	s.db = db
	s.dbMutex.Unlock()
	s.cache.Flush()
}

func (s *server) getLocation(ip net.IP, format string) ([]byte, error) {
	cacheKey := ip.String() + format
	rec, found := s.cache.Get(cacheKey)
	if found {
		log("serving %s from cache", cacheKey)
		return rec.([]byte), nil
	}
	var record interface{}
	var err error

	s.dbMutex.RLock()
	if s.db == nil {
		s.dbMutex.RUnlock()
		return nil, errors.New("not ready")
	}
	if format == "short" {
		r := location{}
		err = s.db.Lookup(ip, &r)
		sr := shortResult{Country: r.Country.ISOCode, City: r.City.Names.EN}
		if len(r.Subdivisions) > 0 {
			sr.Sub = r.Subdivisions[0].IsoCode
		}
		record = sr
	} else {
		r := geoip2.City{}
		err = s.db.Lookup(ip, &r)
		record = r
	}
	s.dbMutex.RUnlock()
	if err != nil {
		log("Lookup: %s", err)
		return nil, err
	}
	resp, err := json.Marshal(record)
	if err != nil {
		log("Marshal: %s", err)
	}
	s.cache.Set(cacheKey, resp, cache.DefaultExpiration)
	return resp, err
}

func (s *server) locationHandler(w http.ResponseWriter, r *http.Request) {
	if e := r.ParseForm(); e != nil {
		log("ParseForm: %s", e)
		http.Error(w, "parse error", http.StatusBadRequest)
		return
	}
	ip := net.ParseIP(r.Form.Get("ip"))
	if ip == nil {
		http.Error(w, "invalid ip", http.StatusBadRequest)
		return
	}

	resp, err := s.getLocation(ip, r.Form.Get("format"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(resp)

}

type location struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	City struct {
		Names struct {
			EN string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"city"`
	Subdivisions []struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"subdivisions"`
}

type shortResult struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Sub     string `json:"sub"`
}
