// Example of creating a web based application purely using
// the net/http package to display weather information and
// Twitter Bootstrap so it doesn't look like it's '92.
//
// To start the app, run:
//    go run weatherweb.go
//
// Accessible via:  http://localhost:8888/here
package main

import (
	"encoding/json"
	"fmt"
	"html/template"

	owm "github.com/briandowns/openweathermap"
	//	"io/ioutil"

	"net/http"
	"os"
)

// URL is a constant that contains where to find the IP locale info
const URL = "http://ip-api.com/json"

// Data will hold the result of the query to get the IP
// address of the caller.
type Data struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	ORG         string  `json:"org"`
	AS          string  `json:"as"`
	Message     string  `json:"message"`
	Query       string  `json:"query"`
}

// getLocation will get the location details for where this
// application has been run from.
func getLocation() (*Data, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	r := &Data{}
	if err = json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}

// getCurrent gets the current weather for the provided location in
// the units provided.
func getCurrent(l, u, lang string) (*owm.CurrentWeatherData, error) {
	w, err := owm.NewCurrent(u, lang, os.Getenv("OWM_API_KEY")) // Create the instance with the given unit
	if err != nil {
		return nil, err
	}
	w.CurrentByName(l) // Get the actual data for the given location
	return w, nil
}

// hereHandler will take are of requests coming in for the "/here" route.
func hereHandler(w http.ResponseWriter, r *http.Request) {
	location, err := getLocation()
	if err != nil {
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}
	wd, err := getCurrent(location.City, "F", "en")
	if err != nil {
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}
	// Process our template
	t, err := template.ParseFiles("templates/here.html")
	if err != nil {
		fmt.Fprint(w, http.StatusInternalServerError)
		return
	}
	// We're doin' naughty things below... Ignoring icon file size and possible errors.
	_, _ = owm.RetrieveIcon("static/img", wd.Weather[0].Icon+".png")

	// Write out the template with the given data
	t.Execute(w, wd)
}

// Run the app
func main() {
	http.HandleFunc("/here", hereHandler)
	// Make sure we can serve our icon files once retrieved
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.ListenAndServe(":8888", nil)
}
