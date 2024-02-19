package utils

import (
	"fmt"
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

var geoData geoip2.Reader

func InitGeo() {
	db, err := geoip2.Open("D:/go/web-service-gin/internal/utils/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	geoData = *db
}

func GetLocationByIp(ip string) {
	ipNet := net.ParseIP(ip)
	record, err := geoData.City(ipNet)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
	fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
	fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
}
