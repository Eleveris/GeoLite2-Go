package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type GeoLite2 interface {
	GetCountry(ip string) (country string, err error)
}

type geoLite2 struct {
	host    string
	account string
	key     string
	client  *http.Client
}

type geoLite2CountryAnswer struct {
	Country struct {
		Iso string `json:"iso_code"`
	} `json:"country"`
}

func (gl geoLite2) GetCountry(ip string) (country string, err error) {
	request, err := http.NewRequest("GET", gl.host+"country/"+ip, nil)
	request.SetBasicAuth(gl.account, gl.key)
	response, err := gl.client.Do(request)
	bodyReader := bufio.NewReader(response.Body)
	rawBody, err := bodyReader.ReadSlice('\r')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	body := geoLite2CountryAnswer{}
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		return "", err
	}
	return body.Country.Iso, nil
}

func NewGeoLite2(host string, account string, key string) GeoLite2 {
	return geoLite2{
		host:    host,
		account: account,
		key:     key,
		client:  &http.Client{},
	}
}
