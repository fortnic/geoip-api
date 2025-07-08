# Fortnic GeoIP API

Fortnic GeoIP is a blazing-fast, open-source IP geolocation API built in Go and optimized for high performance. It is completely free, no account or API key required, no tracking, no user data collection and supports multiple response formats.

> 🚀 Live Demo: [https://geoip.fortnic.com](https://geoip.fortnic.com)

---

## 📆 Features

* Ultra-fast using `gnet` + `fasthttp`
* Redis caching for performance
* IP geolocation lookup with city-level accuracy
* Supports JSON, XML, CSV, JSONP, and plain text
* Fully CORS-enabled (for browser APIs)
* Public endpoint with zero restrictions
* Static file serving for landing page

---

## 📃 Tech Stack

| Technology      | Purpose                            |
| --------------- | ---------------------------------- |
| **Go**          | Core language for the API          |
| **Gnet**        | Non-blocking TCP network library   |
| **Fasthttp**    | Fast HTTP implementation for Go    |
| **Redis**       | Caching layer for IP lookups       |
| **IP2Location** | LITE BIN for geolocation data      |
| **DB-IP**       | MMDB format for ASN database       |

---

## 📊 API Usage

| Endpoint                      | Description                           | 
| ----------------------------- | ----------------------------------    |
| `/ip`                         | Returns visitor's IP (plain text)     | 
| `/ip?format=json\|csv\|xml`   | IP in selected format                 |
| `/me`                         | Visitor IP lookup with geolocation    | 
| `/{ip}`                       | Geolocation lookup for provided IP    | 
| `/me?format=json\|csv\|xml`   | Visitor IP in formatted geolocation   |
| `/{ip}?format=json\|csv\|xml` | Geolocation for IP in chosen format   |

**Supports JSONP**: Use `?callback=yourFunc` for JSONP response.

**Example**:

```bash
curl https://geoip.fortnic.com/1.1.1.1?format=json
```

---

## 📊 Data Fields

* `ip`
* `country_code`
* `country`
* `region`
* `city`
* `postal_code`
* `latitude`
* `longitude`
* `organization` (ASN + ISP)
* `timezone`

---

## 🗂️ Data Sources

Download the latest database files manually:

* **IP2Location DB11 LITE (BIN)**: [https://lite.ip2location.com/ip2location-lite](https://lite.ip2location.com/ip2location-lite)
* **DB-IP ASN (MMDB)**: [https://db-ip.com/db/download/ip-to-asn-lite](https://db-ip.com/db/download/ip-to-asn-lite)

Put the `.BIN` and `.MMDB` files in the `./data` directory.

> You can combine both for a more complete result.

---

## ⚙️ Installation & Build

### Development Mode

```bash
git clone https://github.com/fortnic/geoip-api.git
cd geoip-api
go run .
```

The application is now running at `http://127.0.0.1:8080`

### Production Build + Compression

```bash
go build -ldflags="-s -w" -o geoip-api
upx -9 --ultra-brute geoip-api
```

### Run as a Service (Linux systemd)

```bash
sudo cp -r geoip-api /opt/
sudo cp system/geoip-api.service /etc/systemd/system/
sudo systemctl daemon-reexec
sudo systemctl daemon-reload
sudo systemctl enable geoip-api
sudo systemctl start geoip-api
```

Skip the steps above by running: `./setup.sh`

---

## 🚀 Deploy with CDN (Cloudflare)

Use Cloudflare in front of this server for:

* Global caching
* Automatic compression (Brotli/gzip)
* DDOS protection
* Optional rate limiting

---

## 🌐 Credits

* [IP2Location LITE](https://lite.ip2location.com/)
* [DB-IP Lite](https://db-ip.com/)
* [Gnet by panjf2000](https://github.com/panjf2000/gnet)
* [Fasthttp by valyala](https://github.com/valyala/fasthttp)

---

Made with ❤️ by [Fortnic](https://fortnic.com)
