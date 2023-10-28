package main

import (
	"encoding/json"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var geoLite2API GeoLite2
var geoLite2LocalAPI GeoLite2Local
var cache Cache

func main() {
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatal(err)
	}
	geoLite2API = NewGeoLite2(os.Getenv("GEOLITE_URI"), os.Getenv("GEOLITE_ACCOUNT"), os.Getenv("GEOLITE_KEY"))
	ttl, err := strconv.ParseInt(os.Getenv("CACHE_TTL"), 10, 32)
	if err != nil {
		ttl = 60
	}
	cache = NewMemcached(os.Getenv("MEMCACHED_HOST"), os.Getenv("MEMCACHED_PORT"), int32(ttl))
	geoLite2LocalAPI, err = NewGeoLite2Local(os.Getenv("GEOLITE_DBFILEPATH"))
	if err == nil {
		defer geoLite2LocalAPI.Close()
	}
	port := os.Getenv("LISTEN_PORT")
	http.HandleFunc("/geoip", getCountry)
	http.HandleFunc("/geoip_local", getCountryLocal)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getCountry(w http.ResponseWriter, r *http.Request) {
	log.Println("got a request for remote")
	ip, err := validateRequest(r)
	if err != nil {
		errorResponse(&w, err, 400)
		return
	}
	cached, err := cache.Get("remote." + ip.String())
	log.Println(err)
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		log.Println("cache error: ", err)
	}
	if err == nil {
		successResponse(&w, ip.String(), cached, false)
		return
	}
	country, err := geoLite2API.GetCountry(ip.String())
	if err != nil {
		errorResponse(&w, err, 400)
		return
	}
	err = cache.Set("remote."+ip.String(), country)
	if err != nil {
		log.Println("cache error: ", err)
	}
	successResponse(&w, ip.String(), country, true)
	return
}

func getCountryLocal(w http.ResponseWriter, r *http.Request) {
	log.Println("got a request for local")
	if geoLite2LocalAPI == nil {
		errorResponse(&w, errors.New("файл базы данных geoip2 не доступен"), 500)
		return
	}
	ip, err := validateRequest(r)
	if err != nil {
		errorResponse(&w, err, 400)
		return
	}
	cached, err := cache.Get("local." + ip.String())
	log.Println(err)
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		log.Println("cache error: ", err)
	}
	if err == nil {
		successResponse(&w, ip.String(), cached, false)
		return
	}
	country, err := geoLite2LocalAPI.GetCountry(ip)
	if err != nil {
		errorResponse(&w, err, 400)
		return
	}
	err = cache.Set("local."+ip.String(), country)
	if err != nil {
		log.Println("cache error: ", err)
	}
	successResponse(&w, ip.String(), country, true)
	return
}

func validateRequest(r *http.Request) (ip net.IP, err error) {
	queryParams := r.URL.Query()
	log.Println(queryParams)
	ipParam := queryParams.Get("ip")
	if ipParam == "" {
		return nil, errors.New("в запросе отсутствует параметр ip")
	}
	if ip = net.ParseIP(ipParam); ip == nil {
		return nil, errors.New("параметр ip не является валидным ip адресом")
	}
	return ip, nil
}

func errorResponse(w *http.ResponseWriter, errMsg error, code int) {
	log.Println("error response")
	response, err := json.Marshal(struct {
		Err     string `json:"error"`
		Errcode int    `json:"errcode"`
	}{errMsg.Error(), code})
	if err != nil {
		log.Println(err)
	}
	log.Println(response)
	(*w).Header().Set("Content-Type", "Application/json")
	(*w).WriteHeader(400)
	_, _ = (*w).Write(response)
}

func successResponse(w *http.ResponseWriter, ip string, country string, fresh bool) {
	log.Println("success response")
	response, _ := json.Marshal(struct {
		Ip      string `json:"ip"`
		Country string `json:"country"`
		Fresh   bool   `json:"fresh"`
	}{ip, country, fresh})
	(*w).Header().Set("Content-Type", "Application/json")
	_, _ = (*w).Write(response)
}
