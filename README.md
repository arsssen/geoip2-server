# geoip2-server
HTTP server which provides API to resolve IP addresses to locations using Maxmind's GeoLite2 City db.

Can automatically get(and update) the database using the provided license key, or use existing database file.

## Configuration

Configuration can be passed as command-line flags or environment variables


| Env  	| Flag  	|   Example	|   Description	|
|---	|---	|---	|---	|
| `LICENSE_KEY`  	| `-L` or `--license-key`   	|  `PxnS34uOcQEtCPAA` 	| MaxMind license key  	|
| `LISTEN_ADDR`  	| `-A` or `--listen-addr`  	| `":8080"`   	| Address:port for http server to listen on. Default: :8080  	|   	|
| `DB_FILE`  	|  `-D` or `--db-file` 	| `/etc/GeoLite2-City.mmdb`  	| Path to existing db file  	|   	|
| `UPDATE_PERIOD`  	|  `-U` or `--update-period` 	| `10h30m`  	| Time period to check for updates/download updted db. Default: 10 days  	|   	|
 

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
      "ja": "フランクフルト・アム・マイン",
      "pt-BR": "Frankfurt am Main",
      "ru": "Франкфурт",
      "zh-CN": "法兰克福"
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
      "ja": "ヨーロッパ",
      "pt-BR": "Europa",
      "ru": "Европа",
      "zh-CN": "欧洲"
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
      "ja": "ドイツ連邦共和国",
      "pt-BR": "Alemanha",
      "ru": "Германия",
      "zh-CN": "德国"
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
      "fr": "États Unis",
      "ja": "アメリカ",
      "pt-BR": "EUA",
      "ru": "США",
      "zh-CN": "美国"
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
        "ru": "Гессен"
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
