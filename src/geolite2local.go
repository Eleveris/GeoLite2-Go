package main

import (
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
)

type GeoLite2Local interface {
	GetCountry(ip net.IP) (country string, err error)
	Close()
}

type geoLite2Local struct {
	dbReader *geoip2.Reader
}

func (gl geoLite2Local) GetCountry(ip net.IP) (country string, err error) {
	result, err := gl.dbReader.Country(ip)
	if err != nil {
		return "", err
	}
	return result.Country.IsoCode, nil
}

func (gl geoLite2Local) Close() {
	gl.dbReader.Close()
}

func NewGeoLite2Local(dbPath string) (GeoLite2Local, error) {
	dbReader, err := geoip2.Open(dbPath)
	if err != nil {
		log.Println("geoip open file error", err)
		return nil, err
	}
	return geoLite2Local{dbReader: dbReader}, nil
}
