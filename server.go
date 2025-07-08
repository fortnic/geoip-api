package main

import (
	"bufio"
	"bytes"
	"log"
	"sync"

	"github.com/panjf2000/gnet"
	"github.com/valyala/fasthttp"
)

var bufPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

type geoipServer struct {
	*gnet.EventServer
	router fasthttp.RequestHandler
}

func (gs *geoipServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("GeoIP API listening on %s", srv.Addr.String())
	return
}

func (gs *geoipServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	req := &fasthttp.Request{}

	if err := req.Read(bufio.NewReader(bytes.NewReader(frame))); err == nil {
		req.Header.Set("X-Real-IP", c.RemoteAddr().String())
		ctx := &fasthttp.RequestCtx{}
		ctx.Init(req, c.RemoteAddr(), nil)
		gs.router(ctx)
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()
		writer := bufio.NewWriter(buf)
		ctx.Response.Write(writer)
		writer.Flush()
		resp := buf.Bytes()
		bufPool.Put(buf)
		return resp, gnet.None
	}

	return []byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 0\r\n\r\n"), gnet.None
}
