package main

// run command: $ go run train_simulator.go train.go helping_functions_and_struct.go  user_interaction_functions.go json_data_functions.go


import (
	"sync"
	"fmt"
)

var mutex = &sync.Mutex{}
var secInHour = 1
var babblerMode bool = false
var wg sync.WaitGroup
var reportDate = getReportDate()


func main() {
	clearReportFile()
	fmt.Println("Please determine thow much secounds is one hour in simulation: ")
	fmt.Scan(&secInHour)

	railwaysSwitches := getRailwaySwitches()
	railways := getRailways(railwaysSwitches)
	platforms := getPlatforms(railways)

	railwaysSwitchesAvailable := initRailwaysSwitchesAvailableMap(railwaysSwitches)
	railwaysAvailable := initRailwaysAvailableMap(railways)

	trains := getTrains(&platforms, &railwaysSwitchesAvailable, &railwaysAvailable)
	wg.Add(len(trains))

	getModeFromUser()

	for i := range trains {
		startTrain(&trains[i])
	}
	handleStillMode(trains)

	wg.Wait()
}

func startTrain(t* train) {

	go func() {

		defer wg.Done()

		t.handleChangeStationRailway( 1)

		for {
			t.printInfoCurrentRailway()

			t.simulateCurrentRailwayTime()

			t.handleStationIfExist()

			t.handleUseRailwaySwitch()

			t.handleChangeRailway(0)

			t.trainHandleReleasingStation() // of course if we are on the station
			t.handleReleasingRailwaySwitch()
			t.handleReleasingRailway()
			t.incrementRailway()
		}
	}()

}
