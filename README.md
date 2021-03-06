# geoip2-server

Provides API to resolve IP addresses to locations using Maxmind's GeoLite2 City db.

Can be used as a library or standalone HTTP server.

Can automatically get(and update) the database using the provided license key or use existing database file.

## Configuration

Configuration can be passed as command-line flags or environment variables


| Env  	| Flag  	|   Example	|   Description	|
|---	|---	|---	|---	|
| `LICENSE_KEY`  	| `-L` or `--license-key`   	|  `PxnS34uOcQEtCPAA` 	| MaxMind license key  	|
| `LISTEN_ADDR`  	| `-A` or `--listen-addr`  	| `":8080"`   	| Address:port for http server to listen on. Default: :8080  	|   	|
| `DB_FILE`  	|  `-D` or `--db-file` 	| `/etc/GeoLite2-City.mmdb`  	| Path to existing db file  	|   	|
| `UPDATE_PERIOD`  	|  `-U` or `--update-period` 	| `10h30m`  	| Time period to check for updates/download updted db. Default: 10 days  	|   	|
| `WORK_DIR`  	|  `-W` or `--work-dir` 	| `/.`  	|  Working directory to store downloaded files. If not specified, temp directory will be used  	|   	|
 

## Running (with docker):


Run on port `80`, download the database using provided license key, check for updates once a day:

```bash
docker run --name my-geoip-db -p80:8080 -e LICENSE_KEY=PxnS34uOcQEtCPAA -e UPDATE_PERIOD=24h arsssen/geoip2-server:latest
```


Run using existing db(assuming `GeoLite2-City.mmdb` file is in `/etc` local directory):

```bash
docker run --name my-geoip-db -e DB_FILE=/db/GeoLite2-City.mmdb -v /etc:/db  arsssen/geoip2-server:latest
```



Run starting with existing db, but periodicaly check for updates and download new versions:

```bash
docker run --name my-geoip-db -e DB_FILE=/db/GeoLite2-City.mmdb -v /etc:/db -e LICENSE_KEY=PxnS34uOcQEtCPAA  arsssen/geoip2-server:latest
```

## Embedding in your app

```bash
go get  github.com/arsssen/geoip2-server
```

See `main.go` for example how to use.


## HTTP API

There's a single API endpoint: `/`
IP should be passed as a GET parameter named `ip`
There's also a `format` parameter specifying whether it should return full information, or only country/city/state(`format=short`).

### Examples:

Request:

`/?ip=8.8.8.8`

Response:
```json
{
  "City": {
    "GeoNameID": 2925533,
    "Names": {
      "de": "Frankfurt am Main",
      "en": "Frankfurt am Main",
      "es": "Francfort",
      "fr": "Francfort-sur-le-Main",
      "ja": "??????????????????????????????????????????",
      "pt-BR": "Frankfurt am Main",
      "ru": "??????????????????",
      "zh-CN": "????????????"
    }
  },
  "Continent": {
    "Code": "EU",
    "GeoNameID": 6255148,
    "Names": {
      "de": "Europa",
      "en": "Europe",
      "es": "Europa",
      "fr": "Europe",
      "ja": "???????????????",
      "pt-BR": "Europa",
      "ru": "????????????",
      "zh-CN": "??????"
    }
  },
  "Country": {
    "GeoNameID": 2921044,
    "IsInEuropeanUnion": true,
    "IsoCode": "DE",
    "Names": {
      "de": "Deutschland",
      "en": "Germany",
      "es": "Alemania",
      "fr": "Allemagne",
      "ja": "????????????????????????",
      "pt-BR": "Alemanha",
      "ru": "????????????????",
      "zh-CN": "??????"
    }
  },
  "Location": {
    "AccuracyRadius": 1000,
    "Latitude": 50.1188,
    "Longitude": 8.6843,
    "MetroCode": 0,
    "TimeZone": "Europe/Berlin"
  },
  "Postal": {
    "Code": "60313"
  },
  "RegisteredCountry": {
    "GeoNameID": 6252001,
    "IsInEuropeanUnion": false,
    "IsoCode": "US",
    "Names": {
      "de": "USA",
      "en": "United States",
      "es": "EE. UU.",
      "fr": "??tats Unis",
      "ja": "????????????",
      "pt-BR": "EUA",
      "ru": "??????",
      "zh-CN": "??????"
    }
  },
  "RepresentedCountry": {
    "GeoNameID": 0,
    "IsInEuropeanUnion": false,
    "IsoCode": "",
    "Names": null,
    "Type": ""
  },
  "Subdivisions": [
    {
      "GeoNameID": 2905330,
      "IsoCode": "HE",
      "Names": {
        "de": "Hessen",
        "en": "Hesse",
        "es": "Hessen",
        "fr": "Hesse",
        "ru": "????????????"
      }
    }
  ],
  "Traits": {
    "IsAnonymousProxy": false,
    "IsSatelliteProvider": false
  }
}
```


Request:
`/?ip=8.8.8.8&format=short`

Response:
```json
{"country":"DE","city":"Frankfurt am Main","sub":"HE"}
```
