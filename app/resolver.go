package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/oschwald/maxminddb-golang"
	"github.com/patrickmn/go-cache"
)

//NewResolver starts the server with given options
func NewResolver(ctx context.Context, args Opts) (*resolver, error) {
	logger := func(format string, v ...interface{}) {}
	if args.Logger != nil {
		logger = args.Logger
	}
	s := newResolver(ctx, logger)
	setDB := func(filename string) {
		logger("opening db file %s", filename)
		db, err := maxminddb.Open(filename)
		if err != nil {
			logger("cannot open db %s: %s", filename, err)
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
			return nil, fmt.Errorf("invalid update period: %s", args.UpdatePeriod)
		}
		go s.startUpdater(ctx, period, args.LicenseKey, setDB)
	}
	if args.DBFile == "" && args.LicenseKey == "" {
		return nil, fmt.Errorf("cannot run without license key or existing db file")
	}

	return &s, nil

}

func newResolver(ctx context.Context, logger func(format string, v ...interface{})) resolver {
	return resolver{
		ctx:     ctx,
		dbMutex: &sync.RWMutex{},
		cache:   cache.New(5*time.Hour, time.Hour),
		logger:  logger,
	}
}

type resolver struct {
	ctx     context.Context
	db      *maxminddb.Reader
	dbMutex *sync.RWMutex
	cache   *cache.Cache
	logger  func(format string, v ...interface{})
}

func (s *resolver) GetLocationJSON(ip net.IP, format string) ([]byte, error) {
	cacheKey := ip.String() + format
	rec, found := s.cache.Get(cacheKey)
	if found {
		s.logger("serving %s from cache", cacheKey)
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
		s.logger("Lookup: %s", err)
		return nil, err
	}
	resp, err := json.Marshal(record)
	if err != nil {
		s.logger("Marshal: %s", err)
	}
	s.cache.Set(cacheKey, resp, cache.DefaultExpiration)
	return resp, err
}
