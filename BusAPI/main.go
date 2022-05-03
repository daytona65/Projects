package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BusTimings struct {
	StopID      string `json:"stop_id"`
	StopName    string `json:"stop_name"`
	BusLineID   int    `json:"busline_id"`
	BusLineName string `json:"busline_name"`
	VehicleID   int    `json:"vehicle_id"`
	ArrivalTime int    `json:"arrival_time"`
}

const stopurl = "https://dummy.uwave.sg/busstop/"
const lineurl = "https://dummy.uwave.sg/busline/"
const samplebustimingurl = "https://dummy.uwave.sg/example/busstop"
const samplebuslineurl = "https://dummy.uwave.sg/example/busline"

var StopID = []string{
	"378204", "383050", "378202", "383049", "382998", "378237", "378233", "378230",
	"378229", "378228", "378227", "382995", "378224", "378226", "383010", "383009",
	"383006", "383004", "378234", "383003", "378222", "383048", "378203", "382999",
	"378225", "383014", "383013", "383011", "377906", "383018", "383015", "378207",
}

var LineID = []string{
	"44478", "44479", "44480", "44481",
}

var bustimings []BusTimings

func main() {
	var stopexists bool = false
	var stopNumber string

	//Request input from user and input validation for user's current bus stop.
	for !stopexists {
		fmt.Println("Bus Stop Number: ")
		fmt.Scan(&stopNumber)
		for _, stops := range StopID {
			if stops == stopNumber {
				stopexists = true
			}
		}
		if !stopexists {
			fmt.Println("Bus Stop does not exist.")
			continue
		}
	}

	//Retrieving stop data
	var stopData = getStop(stopurl + stopNumber)
	var busID []int
	var busRoutes = make(map[int][][2]int)

	//For testing. Using sample bus stop URL.================================================================
	myJsonString := `{
						"external_id": "SNDS://27211", 
						"forecast": [{
										"forecast_seconds": 269.1422242217484, 
										"route": {
													"id": 11815, 
													"name": "Bus Campus Loop Red (CL-R) [Singapore]", 
													"short_name": "Campus Loop Red (CL-R)"
												}, 
										"rv_id": 44478, 
										"total_pass": 1432.1273333333334, 
										"vehicle": "31410 PC 3068 E/ at NTU Bus park of NTU", 
										"vehicle_id": 31410
									}], 
						"geometry": [{
										"external_id": null, 
										"lat": "1.3476696531", 
										"lon": "103.6804890633", 
										"seq": 1
									}], 
						"id": 378224, 
						"name": "LWN Library, Opp. NIE", 
						"name_en": null, 
						"name_ru": null, 
						"nameslug": "", 
						"resource_uri": "/routes/api/platformbusarrival/378224/"
					}`
	json.Unmarshal([]byte(myJsonString), &stopData)
	//=========================================================================================================
	if len(stopData.Forecast) == 0 {
		fmt.Println("No buses available.")
		return
	}

	for _, routes := range stopData.Forecast {
		var bus [2]int
		bus[0] = routes.VehicleID
		bus[1] = int(routes.ForecastSeconds) / 60
		busRoutes[routes.RvID] = append(busRoutes[routes.RvID], bus)
		busID = append(busID, routes.VehicleID)
		bustimings = append(bustimings, BusTimings{StopID: strconv.Itoa(stopData.ID), StopName: stopData.Name, BusLineID: routes.Route.ID, BusLineName: routes.Route.Name, VehicleID: routes.VehicleID, ArrivalTime: int(routes.ForecastSeconds) / 60})
		fmt.Println("Bus", strconv.Itoa(bus[0]), "on", routes.Route.ShortName, "is arriving in", bus[1], "minutes.")
	}
}

func GetTimingEndpoint(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(bustimings)
}

func CreateTimingEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var busdata BusTimings
	_ = json.NewDecoder(r.Body).Decode(&bustimings)
	busdata.StopID = params["stop_id"]
	busdata.StopName = params["stop_name"]
	busdata.BusLineID, _ = strconv.Atoi(params["busline_id"])
	busdata.BusLineName = params["busline_name"]
	busdata.VehicleID, _ = strconv.Atoi(params["vehicle_id"])
	busdata.ArrivalTime, _ = strconv.Atoi(params["arrival_time"])
	bustimings = append(bustimings, busdata)
	json.NewEncoder(w).Encode(bustimings)
}

func DeleteTimingEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range bustimings {
		var a int
		a, _ = strconv.Atoi(params["vehicle_id"])
		if item.StopID == params["stop_id"] && item.VehicleID == a {
			bustimings = append(bustimings[:index], bustimings[:index+1]...)
		}
	}
	json.NewEncoder(w).Encode(bustimings)
}

func router() {
	router := mux.NewRouter()
	router.HandleFunc("/bustimings", GetTimingEndpoint).Methods("GET")
	router.HandleFunc("/bustimings", CreateTimingEndpoint).Methods("POST")
	router.HandleFunc("/bustimings", DeleteTimingEndpoint).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":12345", router))
}
