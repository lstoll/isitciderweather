package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type GeoResp struct {
	Latitude  float32
	Longitude float32
}

type WeatherResp struct {
	WeatherObservation WeatherObservation
}

type WeatherObservation struct {
	Temperature string
}

type PageContent struct {
	ResponseText string
}

var chttp = http.NewServeMux()

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	chttp.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/", HomeHandler) // homepage
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.Path, ".") {
		chttp.ServeHTTP(w, r)
	} else {
		var responseText string
		// fmt.Fprintf(w, "HomeHandler")
		xff := r.Header.Get("X-Forwarded-For")
		var clientIp string
		if xff == "" {
			clientIp = r.RemoteAddr
		} else {
			clientIp = xff
		}
		loc, err := LocationForIP(clientIp)
		if err != nil {
			log.Println("Error locating IP: " + clientIp)
			responseText = "Couldn't find you!"
		}
		log.Println(loc)
		weather := WeatherForLoc(loc)
		log.Println(weather.Temperature)
		temp, _ := strconv.ParseFloat(weather.Temperature, 32)
		if temp > 19 {
			responseText = "Temp is more than 19!"
		} else {
			responseText = "Temp is less than 19"
		}

		t := template.New("index.html")     //create a new template
		t, err = t.ParseFiles("index.html") //open and parse a template text file
		if err != nil {
			log.Fatalln(err)
		}
		err = t.Execute(w, PageContent{responseText}) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func WeatherForLoc(loc GeoResp) WeatherObservation {
	url := fmt.Sprintf("http://api.geonames.org/findNearByWeatherJSON?lat=%g&lng=%g&username=%s", loc.Latitude, loc.Longitude, os.Getenv("GEONAMES_USERNAME"))
	log.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	weather := WeatherResp{}
	err = decoder.Decode(&weather)
	if err != nil {
		log.Fatalln(err)
	}
	return weather.WeatherObservation
}

func LocationForIP(ip string) (GeoResp, error) {
	// Get the location
	url := "http://freegeoip.net/json/" + "178.19.210.162" // + ip
	log.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return GeoResp{}, errors.New("Not found")
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	log.Println(resp.Body)
	loc := GeoResp{}
	err = decoder.Decode(&loc)
	if err != nil {
		log.Fatalln(err)
	}
	return loc, nil
}
