package main

import (
	"log"

	"github.com/ip2location/ip2location-go"
	"github.com/oschwald/geoip2-golang"
)

var (
	ip2db *ip2location.DB
	asnDB *geoip2.Reader
)

func init() {
	var err error
	ip2db, err = ip2location.OpenDB("data/IP2LOCATION-LITE-DB11.IPV6.BIN")
	if err != nil {
		log.Fatal("Failed to open:", err)
	}

	asnDB, err = geoip2.Open("data/dbip-asn-lite-2025-07.mmdb")
	if err != nil {
		log.Fatal("Failed to open:", err)
	}
}
