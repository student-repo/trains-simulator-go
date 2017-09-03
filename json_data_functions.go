package main

import (
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
	"time"
)

type JSONTrain struct {
	Name     string `json:"name"`
	Route    []int `json:"route"`
	MaxSpeed int `json:"max_speed"`
	Color    string `json:"color"`
}

type JSONDataStructure struct {
	Trains          []JSONTrain `json:"trains"`
	Railways        []railwayJSON `json:"railways"`
	RailwaySwitches []railwaySwitch `json:"railway_switches"`
	ReportYear      int `json:"report_year"`
	ReportMonth      int `json:"report_month"`
	ReportDay      int `json:"report_day"`
	ReportHour      int `json:"report_hour"`
	ReportMinute      int `json:"report_minutes"`
	ReportSecond      int `json:"report_seconds"`

}


type railwayJSON struct {
	MaxSpeed       int `json:"max_speed"`
	Length         int `json:"length"`
	RailwaySwitch1ID int `json:"railway_switch_1_ID"`
	RailwaySwitch2ID int `json:"railway_switch_2_ID"`
	Platforms        []int `json:"platforms"`
	Id             int `json:"id"`
}

func getRailwaySwitches() []railwaySwitch {
	k := getJSONData().RailwaySwitches
	s := make([]railwaySwitch,len(k))
	for i, k := range k {
		s[i] = railwaySwitch{k.UseTime, k.Id}
	}
	return s
}

func (r railway) toString() string {
	return toJson(r)
}

func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

func getPlatforms(railways map[int]railway)map[int]map[int]platform {
	platforms := map[int]map[int]platform{}
	for _,k := range railways {
		if len(k.platforms) != 0 {
			platforms[k.Id] = k.platforms
		}
	}
	return platforms
}

func getTrains(platforms* map[int]map[int]platform, railwaySwitchAvailable* map[int]bool, railwayAvailable* map[int]bool) []train {
	r := getRailways(getRailwaySwitches())
	t := getJSONData().Trains
	tt := make([]train,len(t))

	for i,k := range t {
		rr := make([]railway,len(k.Route))
		for ii, kk := range k.Route {
			rr[ii] = r[kk]
		}
		tt[i] = train{k.Name, rr, k.MaxSpeed, 0, k.Color,
			*platforms, *railwaySwitchAvailable, *railwayAvailable, 0, 0}
	}
	return tt
}


func getReportDate() time.Time {

		return time.Date(
			getJSONData().ReportYear, time.Month(getJSONData().ReportMonth), getJSONData().ReportDay, getJSONData().ReportHour,
			getJSONData().ReportMinute, getJSONData().ReportSecond, 651387237, time.UTC)

}


func getJSONData() JSONDataStructure {
	raw, err := ioutil.ReadFile("../src/config2.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var c JSONDataStructure
	json.Unmarshal(raw, &c)
	return c
}

func getRailways(railwaySwitches []railwaySwitch)map[int]railway {
	railwaysJSON := getJSONData().Railways
	//s := make([]railway,len(railwaysJSON))
	s := make(map[int]railway)
	l := 0
	for i, key := range railwaysJSON {
		for _, k := range railwaySwitches {
			if (k.Id == key.RailwaySwitch1ID || k.Id == key.RailwaySwitch2ID) && l == 0 {
				asd := make(map[int]platform, len(key.Platforms))
				for iii, kkk := range key.Platforms {
					asd[iii] = platform{true, kkk, iii}
				}
				s[key.Id] = railway{key.MaxSpeed, key.Length, k,k,
						    asd, -1, key.Id}
				l++
			} else if (k.Id == key.RailwaySwitch1ID || k.Id == key.RailwaySwitch2ID) && l == 1 {
				s[key.Id] = railway{s[i].MaxSpeed, s[i].Length, s[i].RailwaySwitch1,
						    k, s[i].platforms, s[i].Platform, s[i].Id}
				break
			}
		}
		l = 0
	}
	return s
}
