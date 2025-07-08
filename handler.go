package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

func ipHandler(ctx *fasthttp.RequestCtx) {
	ip := getRealIP(ctx)
	format := string(ctx.QueryArgs().Peek("format"))
	callback := string(ctx.QueryArgs().Peek("callback"))

	switch format {
	case "json":
		if callback != "" {
			ctx.SetContentType("application/javascript")
			fmt.Fprintf(ctx, "%s({\"ip\": \"%s\"});", callback, ip)
		} else {
			ctx.SetContentType("application/json")
			fmt.Fprintf(ctx, `{"ip": "%s"}`, ip)
		}
	case "xml":
		ctx.SetContentType("application/xml")
		fmt.Fprintf(ctx, `<IP>%s</IP>`, ip)
	case "csv":
		ctx.SetContentType("text/csv")
		fmt.Fprintf(ctx, "ip\n%s", ip)
	default:
		ctx.SetContentType("text/plain")
		fmt.Fprintf(ctx, "%s", ip)
	}
}

func geoHandler(ctx *fasthttp.RequestCtx) {
	path := strings.TrimPrefix(string(ctx.Path()), "/")
	if path == "me" {
		path = getRealIP(ctx)
	}

	ip := net.ParseIP(path)
	if ip == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Invalid IP")
		return
	}

	format := string(ctx.QueryArgs().Peek("format"))
	callback := string(ctx.QueryArgs().Peek("callback"))

	cacheKey := "geoip:" + path
	cached, err := rdb.Get(rctx, cacheKey).Result()
	if err == nil && format == "json" && callback == "" {
		ctx.SetContentType("application/json")
		ctx.WriteString(cached)
		return
	}

	if err == nil {
		var result GeoIPResult
		if jsonErr := json.Unmarshal([]byte(cached), &result); jsonErr == nil {
			response(ctx, result, format, callback)
			return
		}
	}

	rec, err := ip2db.Get_all(path)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Failed to lookup IP")
		return
	}

	asn, err := asnDB.ASN(ip)
	org := ""
	if err == nil {
		org = fmt.Sprintf("AS%d %s", asn.AutonomousSystemNumber, asn.AutonomousSystemOrganization)
	}

	result := GeoIPResult{
		IP:           path,
		CountryCode:  rec.Country_short,
		Country:      rec.Country_long,
		Region:       rec.Region,
		City:         rec.City,
		PostalCode:   rec.Zipcode,
		Latitude:     float64(rec.Latitude),
		Longitude:    float64(rec.Longitude),
		Timezone:     rec.Timezone,
		Organization: org,
	}

	if jsonData, err := json.Marshal(result); err == nil {
		rdb.Set(rctx, cacheKey, jsonData, 24*time.Hour)
	}

	response(ctx, result, format, callback)
}

func fsHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())

	if path == "/" {
		path = "/index.html"
	}

	cleaned := filepath.Clean(path)
	fullPath := filepath.Join("./public", cleaned)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetContentType("text/plain")
		ctx.SetBodyString("Not Found")
		return
	}

	fasthttp.ServeFile(ctx, fullPath)
}
