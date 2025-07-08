package main

import (
	"log"
	"net"
	"strings"

	"github.com/panjf2000/gnet"
	"github.com/valyala/fasthttp"
)

func main() {
	router := func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("X-Powered-By", "Fortnic GeoIP API v1.0; https://geoip.fortnic.com")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type")

		method := string(ctx.Method())

		switch method {
		case fasthttp.MethodGet, fasthttp.MethodHead, fasthttp.MethodOptions:
		default:
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			ctx.SetContentType("text/plain")
			ctx.SetBodyString("405 Method Not Allowed")
			return
		}

		if method == fasthttp.MethodOptions {
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}

		isHead := method == fasthttp.MethodHead

		path := string(ctx.Path())

		switch {
		case path == "/ip":
			ipHandler(ctx)
		case path == "/me":
			geoHandler(ctx)
		default:
			if net.ParseIP(strings.TrimPrefix(path, "/")) != nil {
				geoHandler(ctx)
			} else {
				fsHandler(ctx)
			}
		}

		if isHead {
			ctx.Response.ResetBody()
		}
	}

	server := &geoipServer{router: router}

	log.Fatal(gnet.Serve(server, "tcp://:8080", gnet.WithMulticore(true)))
}
