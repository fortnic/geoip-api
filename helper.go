package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

func writeJSON(ctx *fasthttp.RequestCtx, data any) {
	b, _ := json.Marshal(data)
	ctx.Write(b)
}

func writeJSONP(ctx *fasthttp.RequestCtx, data any, callback string) {
	b, _ := json.Marshal(data)
	fmt.Fprintf(ctx, "%s(%s);", callback, b)
}

func writeXML(ctx *fasthttp.RequestCtx, data GeoIPResult) {
	b, _ := xml.MarshalIndent(data, "", "  ")
	ctx.Write(b)
}

func writeCSV(ctx *fasthttp.RequestCtx, data GeoIPResult) {
	header := "ip,country_code,country,region,city,postal_code,latitude,longitude,organization,timezone\n"
	body := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%f,%f,%s,%s\n",
		data.IP, data.CountryCode, data.Country, data.Region,
		data.City, data.PostalCode, data.Latitude, data.Longitude,
		data.Organization, data.Timezone,
	)
	ctx.WriteString(header + body)
}

func writePlain(ctx *fasthttp.RequestCtx, data GeoIPResult) {
	body := fmt.Sprintf("IP: %s\nCountry Code: %s\nCountry: %s\nRegion: %s\nCity: %s\nPostal Code: %s\nLatitude: %f\nLongitude: %f\nOrganization: %s\nTimezone: %s\n",
		data.IP, data.CountryCode, data.Country, data.Region,
		data.City, data.PostalCode, data.Latitude, data.Longitude,
		data.Organization, data.Timezone,
	)
	ctx.WriteString(body)
}

func getRealIP(ctx *fasthttp.RequestCtx) string {
	if ip := ctx.Request.Header.Peek("CF-Connecting-IP"); len(ip) > 0 {
		return string(ip)
	}

	if ip := ctx.Request.Header.Peek("X-Forwarded-For"); len(ip) > 0 {
		parts := strings.Split(string(ip), ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	return string(ctx.Request.Header.Peek("X-Real-IP"))
}

func response(ctx *fasthttp.RequestCtx, result GeoIPResult, format, callback string) {
	switch format {
	case "json":
		if callback != "" {
			ctx.SetContentType("application/javascript")
			writeJSONP(ctx, result, callback)
		} else {
			ctx.SetContentType("application/json")
			writeJSON(ctx, result)
		}
	case "xml":
		ctx.SetContentType("application/xml")
		writeXML(ctx, result)
	case "csv":
		preview := string(ctx.QueryArgs().Peek("preview")) == "true"
		if preview {
			ctx.SetContentType("text/plain")
		} else {
			ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=geoip-%s.csv", result.IP))
			ctx.SetContentType("text/csv")
		}
		writeCSV(ctx, result)
	default:
		ctx.SetContentType("text/plain")
		writePlain(ctx, result)
	}
}
