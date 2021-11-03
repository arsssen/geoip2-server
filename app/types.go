package app

//Opts is the set of the options needed to run the resolver
type Opts struct {
	LicenseKey   string
	DBFile       string
	UpdatePeriod string
	Logger       func(format string, v ...interface{})
	WorkingDir   string
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
