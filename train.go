package main

import (
	"github.com/mitchellh/colorstring"
	"time"
	"strconv"
)

type train struct {
	name string
	route []railway
	maxSpeed int
	currentRailwayIndex int
	color string
	platforms map[int]map[int]platform
	railwaySwitchAvailable map[int]bool
	railwaysAvailable map[int]bool
	timeFromStart int
	yearsFromStart int
}

func (t train) getNextRailwaySwitch() railwaySwitch{
	if t.currentRailwayIndex < len(t.route) - 1 {
		if t.route[t.currentRailwayIndex].RailwaySwitch1.Id == t.route[t.currentRailwayIndex + 1].RailwaySwitch2.Id ||
			t.route[t.currentRailwayIndex].RailwaySwitch1.Id == t.route[t.currentRailwayIndex + 1].RailwaySwitch1.Id {
			return t.route[t.currentRailwayIndex].RailwaySwitch1
		}
		return t.route[t.currentRailwayIndex].RailwaySwitch2
	} else {
		if t.route[t.currentRailwayIndex].RailwaySwitch1.Id == t.route[0].RailwaySwitch2.Id || t.route[t.currentRailwayIndex].RailwaySwitch1.Id == t.route[0].RailwaySwitch1.Id {
			return t.route[t.currentRailwayIndex].RailwaySwitch1
		}
		return t.route[t.currentRailwayIndex].RailwaySwitch2
	}
}

func (t train) getPreviousRailwaySwitch() railwaySwitch{
	if t.currentRailwayIndex != 0 {
		if t.route[t.currentRailwayIndex].RailwaySwitch1.Id == t.route[t.currentRailwayIndex - 1].RailwaySwitch2.Id ||
			t.route[t.currentRailwayIndex].RailwaySwitch1.Id == t.route[t.currentRailwayIndex - 1].RailwaySwitch1.Id {
			return t.route[t.currentRailwayIndex].RailwaySwitch1
		}
		return t.route[t.currentRailwayIndex].RailwaySwitch2
	} else {
		if t.route[t.currentRailwayIndex].RailwaySwitch1.Id == t.route[len(t.route) - 1].RailwaySwitch2.Id ||
			t.route[t.currentRailwayIndex].RailwaySwitch1.Id == t.route[len(t.route) - 1].RailwaySwitch1.Id {
			return t.route[t.currentRailwayIndex].RailwaySwitch1
		}
		return t.route[t.currentRailwayIndex].RailwaySwitch2
	}
}

func (t train) getCurrentSpeed()int {
	if t.maxSpeed < t.route[t.currentRailwayIndex].MaxSpeed {
		return t.maxSpeed
	}
	return t.route[t.currentRailwayIndex].MaxSpeed
}

func (t* train) incrementRailway() {
	if t.currentRailwayIndex < len(t.route) - 1 {
		t.currentRailwayIndex++
	} else {
		t.currentRailwayIndex = 0
	}
}

func (t train) getNextRailwayIndex()int {
	if t.currentRailwayIndex < len(t.route) - 1 {
		return t.currentRailwayIndex + 1
	}
	return 0
}

func (t train) getRailwayTravelTime()int {
	return int(float64(t.route[t.currentRailwayIndex].Length ) / float64(float64(float64(t.getCurrentSpeed()) * 1000.0 / 3600.0)))
}


func (t* train) handleUseRailwaySwitch() {

	for {
		mutex.Lock()
		railwaySwitchIsFree := t.railwaySwitchAvailable[t.getNextRailwaySwitch().Id]
		mutex.Unlock()

		if railwaySwitchIsFree {
			mutex.Lock()
			t.railwaySwitchAvailable[t.getNextRailwaySwitch().Id] = false
			mutex.Unlock()
			printIfBabblerMode(addColorToString(t.color, sumStrings("Train", t.name, "is on the railway switch with id:",
				strconv.Itoa(t.getNextRailwaySwitch().Id), "we have to wait:", strconv.Itoa(t.getNextRailwaySwitch().UseTime),
				"(", strconv.FormatFloat(getTimeToWaitSec(t.getNextRailwaySwitch().UseTime), 'E', -1, 64), ")")))
			t.appendTimeFromStart(t.getNextRailwaySwitch().UseTime)
			time.Sleep(time.Millisecond * time.Duration(getTimeToWaitInMilliseconds(t.getNextRailwaySwitch().UseTime)))
			break
		}
		printIfBabblerMode(addColorToString(t.color, sumStrings("Train", t.name, "is before railway switch with id:",
			strconv.Itoa(t.getNextRailwaySwitch().Id), "we have to wait until it will be free:")))
		t.appendTimeFromStart(5)
		time.Sleep(time.Millisecond * time.Duration(getTimeToWaitInMilliseconds(5)))
	}
}

func (t* train) handleChangeRailway(start int) {
	if len(t.route[t.getNextRailwayIndex()].platforms) != 0 {
		t.handleChangeStationRailway(0)
	} else {
		t.handleChangeNormalRailway(start)
	}
}

func (t* train) appendTimeFromStart(time int) {
	if t.timeFromStart > 31536000 {
		t.timeFromStart = t.timeFromStart % 31536000
		t.yearsFromStart ++
	}
	t.timeFromStart += time
}

func (t train) printInfoCurrentRailway() {
	printIfBabblerMode(addColorToString(t.color, sumStrings("Train", t.name,
		"is going from railway swich with id:", strconv.Itoa(t.getPreviousRailwaySwitch().Id), "to",
		strconv.Itoa(t.getNextRailwaySwitch().Id), "travel time:", strconv.Itoa(t.getRailwayTravelTime()),
		"(", strconv.FormatFloat(getTimeToWaitSec(t.getRailwayTravelTime()), 'E', -1, 64), ")")))
}

func (t* train) handleReleasingRailwaySwitch() {
	id := t.getNextRailwaySwitch().Id
	printIfBabblerMode(addColorToString(t.color, sumStrings(t.name, "left switch railway with id:", strconv.Itoa(id))))
	mutex.Lock()
	t.railwaySwitchAvailable[id] = true
	mutex.Unlock()
}

func (t* train) handleReleasingRailway() {
	id := t.route[t.currentRailwayIndex].Id
	printIfBabblerMode(addColorToString(t.color, sumStrings(t.name, "left railway with id:", strconv.Itoa(id))))
	mutex.Lock()
	t.railwaysAvailable[id] = true
	mutex.Unlock()
}

func (t* train) simulateCurrentRailwayTime() {
	time.Sleep(time.Millisecond * time.Duration(getTimeToWaitInMilliseconds(t.getRailwayTravelTime())))
	t.appendTimeFromStart(t.getRailwayTravelTime())
}


func (t* train) handleChangeNormalRailway( start int) {
	id := t.route[t.getNextRailwayIndex() - start].Id
	for {
		mutex.Lock()
		railwayIsFree := t.railwaysAvailable[id]
		mutex.Unlock()
		if railwayIsFree {
			mutex.Lock()
			t.railwaysAvailable[id] = false
			mutex.Unlock()
			break
		}
		printIfBabblerMode(addColorToString(t.color, sumStrings("Train", t.name,
			"have to wait before railway with id: ", strconv.Itoa(id), "until it will be free")))
		t.appendTimeFromStart(5)
		time.Sleep(time.Millisecond * time.Duration(getTimeToWaitInMilliseconds(5)))
	}
}

func (t* train) handleChangeStationRailway( start int) {

	for {
		platformID, freePlatformExist := findFirstAvailablePlatform(mutex, t.platforms, t.route[t.getNextRailwayIndex() - start].Id)

		if freePlatformExist {
			mutex.Lock()
			t.platforms[t.route[t.getNextRailwayIndex() - start].Id][platformID] =
				platform{false, t.platforms[t.route[t.getNextRailwayIndex() - start].Id][platformID].downtime, platformID}
			mutex.Unlock()
			printIfBabblerMode(addColorToString(t.color, sumStrings("Train", t.name, "has reserved platform with platformID:", strconv.Itoa(platformID),
				"on station with id:", strconv.Itoa(t.route[t.getNextRailwayIndex() - start].Id))))
			t.route[t.getNextRailwayIndex() - start].Platform = platformID
			break

		}
		printIfBabblerMode(addColorToString(t.color, sumStrings("Train", t.name, "have to wait for free platform, on station with id: ",
			strconv.Itoa(t.route[t.getNextRailwayIndex() - start].Id))))
		t.appendTimeFromStart(5)
		time.Sleep(time.Millisecond * time.Duration(getTimeToWaitInMilliseconds(5)))
	}
}


func (t* train) handleStationIfExist() {
	if len(t.route[t.currentRailwayIndex].platforms) != 0 {
		mutex.Lock()
		report := sumStrings("Train", t.name, "arrived on station with id:",
			strconv.Itoa(t.route[t.currentRailwayIndex].Id), "on a platform",
			strconv.Itoa(t.route[t.currentRailwayIndex].Platform), "time: ", getStationArriveTime(t).String(),"\n")
		writeReportToFile(report)

		mutex.Unlock()
		printIfBabblerMode(addColorToString(t.color, sumStrings("Train", t.name, "is staying in station with id:",
			strconv.Itoa(t.route[t.currentRailwayIndex].Id), "on a platform",
			strconv.Itoa(t.route[t.currentRailwayIndex].Platform), "time:",
			strconv.Itoa(t.route[t.currentRailwayIndex].platforms[t.route[t.currentRailwayIndex].Platform].downtime),
			"(", strconv.FormatFloat(getTimeToWaitSec(t.getNextRailwaySwitch().UseTime), 'E', -1, 64), ")")))
		t.appendTimeFromStart(t.route[t.currentRailwayIndex].platforms[t.route[t.currentRailwayIndex].Platform].downtime)
		time.Sleep(time.Millisecond * time.Duration(getTimeToWaitInMilliseconds(t.route[t.currentRailwayIndex].platforms[t.route[t.currentRailwayIndex].Platform].downtime)))
	}
}

func (t* train)trainHandleReleasingStation() {
	if len(t.route[t.currentRailwayIndex].platforms) != 0 {
		mutex.Lock()
		t.platforms[t.route[t.currentRailwayIndex].Id][t.route[t.currentRailwayIndex].Platform] =
			platform{true, t.platforms[t.route[t.currentRailwayIndex].Id][t.route[t.currentRailwayIndex].Platform].downtime,
				t.route[t.currentRailwayIndex].Platform}
		mutex.Unlock()
		printIfBabblerMode(addColorToString(t.color, sumStrings(t.name, "left station with id:",
			strconv.Itoa(t.route[t.currentRailwayIndex].Id), "and platform id", strconv.Itoa(t.route[t.currentRailwayIndex].Platform))))
		t.route[t.currentRailwayIndex].Platform = -1

	}
}

func printIfBabblerMode(info string){
	if babblerMode {
		colorstring.Println(info)
	}
}

