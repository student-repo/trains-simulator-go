package main

import (
	"bytes"
	"sync"
	"time"
)

type railwaySwitch struct {
	UseTime int `json:"use_time"`
	Id int `json:"id"`
}

type platform struct {
	free bool
	downtime int
	id int
}



type railway struct {
	MaxSpeed int
	Length int
	RailwaySwitch1 railwaySwitch
	RailwaySwitch2 railwaySwitch
	platforms map[int]platform
	Platform int
	Id int
}

type platformInfo struct {
	stationID int
	platformID int
}

func initRailwaysSwitchesAvailableMap(arr []railwaySwitch)map[int]bool{
	railwaysSwitchesAvailable := make(map[int]bool)
	for _, rs := range arr {
		railwaysSwitchesAvailable[rs.Id] = true
	}
	return railwaysSwitchesAvailable
}

func initRailwaysAvailableMap(arr map[int]railway)map[int]bool{
	railwaysAvailable := make(map[int]bool)
	for _, rs := range arr {
		railwaysAvailable[rs.Id] = true
	}
	return railwaysAvailable
}

func addColorToString (color string, s string) string{
	var buffer bytes.Buffer

	buffer.WriteString("[")
	buffer.WriteString(color)
	buffer.WriteString("]")
	buffer.WriteString(s)
	return buffer.String()
}

func sumStrings(nums ...string)string {

	var buffer bytes.Buffer
	for _, s := range nums {
		buffer.WriteString(s)
		buffer.WriteString(" ")
	}
	return buffer.String()
}

func findFirstAvailablePlatform (mutex *sync.Mutex,platforms map[int]map[int]platform, id int)(int, bool) {
	platformID := 0
	platformIsFree := false
	mutex.Lock()
	for _, s := range platforms[id] {
		if s.free {
			platformID = s.id
			platformIsFree = true
			break
		}
	}
	mutex.Unlock()
	return platformID, platformIsFree
}

func getStationArriveTime(t* train) time.Time {
	return reportDate.Add(time.Second * time.Duration(t.timeFromStart)).AddDate(t.yearsFromStart,0,0)
}




