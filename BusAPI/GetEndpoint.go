package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Stop struct {
	ExternalID string `json:"external_id"`

	//Forecast for bus routes and buses
	Forecast []struct {
		ForecastSeconds float64 `json:"forecast_seconds"`
		Route           struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			ShortName string `json:"short_name"`
		} `json:"route"`
		RvID      int     `json:"rv_id"`
		TotalPass float64 `json:"total_pass"`
		Vehicle   string  `json:"vehicle"`
		VehicleID int     `json:"vehicle_id"`
	} `json:"forecast"`

	//Location of the Stop
	Geometry []struct {
		ExternalID interface{} `json:"external_id"`
		Lat        string      `json:"lat"`
		Lon        string      `json:"lon"`
		Seq        int         `json:"seq"`
	} `json:"geometry"`

	ID          int         `json:"id"`
	Name        string      `json:"name"`
	NameEn      interface{} `json:"name_en"`
	NameRu      interface{} `json:"name_ru"`
	Nameslug    string      `json:"nameslug"`
	ResourceURI string      `json:"resource_uri"`
}

type Line struct {
	ExternalID  interface{} `json:"external_id"`
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	NameEn      interface{} `json:"name_en"`
	NameRu      interface{} `json:"name_ru"`
	Nameslug    interface{} `json:"nameslug"`
	ResourceURI string      `json:"resource_uri"`
	Routename   string      `json:"routename"`
	Vehicles    []struct {
		Bearing    int    `json:"bearing"`
		DeviceTs   string `json:"device_ts"`
		Enterprise struct {
			EnterpriseID   int    `json:"enterprise_id"`
			EnterpriseName string `json:"enterprise_name"`
		} `json:"enterprise"`
		Lat  string `json:"lat"`
		Lon  string `json:"lon"`
		Park struct {
			ParkID   int    `json:"park_id"`
			ParkName string `json:"park_name"`
		} `json:"park"`
		Position struct {
			Bearing  int    `json:"bearing"`
			DeviceTs int    `json:"device_ts"`
			Lat      string `json:"lat"`
			Lon      string `json:"lon"`
			Speed    int    `json:"speed"`
			Ts       int    `json:"ts"`
		} `json:"position"`
		Projection struct {
			EdgeDistance    string `json:"edge_distance"`
			EdgeID          int    `json:"edge_id"`
			EdgeProjection  string `json:"edge_projection"`
			EdgeStartNodeID int    `json:"edge_start_node_id"`
			EdgeStopNodeID  int    `json:"edge_stop_node_id"`
			Lat             string `json:"lat"`
			Lon             string `json:"lon"`
			OrigLat         string `json:"orig_lat"`
			OrigLon         string `json:"orig_lon"`
			RoutevariantID  int    `json:"routevariant_id"`
			Ts              int    `json:"ts"`
		} `json:"projection"`
		RegistrationCode string `json:"registration_code"`
		RoutevariantID   int    `json:"routevariant_id"`
		Speed            string `json:"speed"`
		Stats            struct {
			AvgSpeed    string `json:"avg_speed"`
			Bearing     int    `json:"bearing"`
			CummSpeed10 string `json:"cumm_speed_10"`
			CummSpeed2  string `json:"cumm_speed_2"`
			DeviceTs    int    `json:"device_ts"`
			Lat         string `json:"lat"`
			Lon         string `json:"lon"`
			Speed       int    `json:"speed"`
			Ts          int    `json:"ts"`
		} `json:"stats"`
		Ts        string `json:"ts"`
		VehicleID int    `json:"vehicle_id"`
	} `json:"vehicles"`
	Via interface{} `json:"via"`
}

func getStop(url string) Stop {
	response, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	var stop Stop
	data, err := ioutil.ReadAll(response.Body)

	if err == nil && data != nil {
		err = json.Unmarshal(data, &stop)
	}

	return stop
}

func getLine(url string) Line {
	response, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	var line Line
	data, err := ioutil.ReadAll(response.Body)

	if err == nil && data != nil {
		err = json.Unmarshal(data, &line)
	}

	return line
}
