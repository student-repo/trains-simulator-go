package main

import (
	"bufio"
	"os"
	"github.com/mitchellh/colorstring"
	"fmt"
	"strconv"
	"io/ioutil"
)

func getModeFromUser(){
	fmt.Println("Please select mode: ")
	fmt.Println("1. Babbler mode")
	fmt.Println("2. Still mode")
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if text == "1\n" {
			babblerMode = true
			break
		} else if text == "2\n" {
			babblerMode = false
			fmt.Println("Still mode is running")
			break
		}
		fmt.Println("Incorrect intput !")
	}
}

func getTimeToWaitInMilliseconds(ttt int) int {
	return (secInHour * ttt * 1000) / 3600
}

func getTimeToWaitSec(ttt int) float64 {
	return (float64(secInHour) * float64(ttt)) / float64(3600)
}

func getTrainFromUser(t []train) train{
	fmt.Println()
	fmt.Println()
	fmt.Println("Please select train: ")
	for i := range t {
		fmt.Println(i, ". ", t[i].name)
	}
	for {
		var i int
		fmt.Scan(&i)
			return t[i]
		colorstring.Println(addColorToString("red", sumStrings("Incorrect intput !")))
	}
}

func printTrainInfo(tt []train) {
	t := getTrainFromUser(tt)
	fmt.Println()
	fmt.Println()
	fmt.Println("Please select: ")
	fmt.Println("1. Get current railway id")
	fmt.Println("2. Get current speed")
	for {
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		if reader.Text() == "1" {
			colorstring.Println(addColorToString("red", sumStrings(t.name, "current railway id: ", strconv.Itoa(t.route[t.currentRailwayIndex].Id))))
			break
		} else if reader.Text() == "2" {
			colorstring.Println(addColorToString("red", sumStrings(t.name, "current speed: ", strconv.Itoa(t.getCurrentSpeed()))))
			break
		}
		colorstring.Println(addColorToString("red", sumStrings("Incorrect intput !")))
	}
}

func handleStillMode(t []train) {
	go func() {
		if !babblerMode {
			for {
				printTrainInfo(t)
			}
		}
	}()
}

func getSecInHour()int64 {

	fmt.Println()
	fmt.Println()
	fmt.Print("Please determine thow much secounds is one hour in simulation: ")
	for {
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		if _, err := strconv.Atoi(reader.Text()); err == nil {
			i, _ := strconv.ParseInt(reader.Text(), 10, 32)
			return i
		}
		colorstring.Println(addColorToString("red", sumStrings("Incorrect intput !")))
		fmt.Print("Please determine thow much secounds is one hour in simulation: ")
	}
}

func getSimulationTime()int64 {

	fmt.Println()
	fmt.Println()
	fmt.Print("Please pass simulation time in seconds: ")
	for {
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		if _, err := strconv.Atoi(reader.Text()); err == nil {
			i, _ := strconv.ParseInt(reader.Text(), 10, 32)
			return i
		}
		colorstring.Println(addColorToString("red", sumStrings("Incorrect intput !")))
		fmt.Print("Please pass simulation time in seconds: ")
	}
}

func writeReportToFile(report1 string) {
	f, err := os.OpenFile("./report", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(report1); err != nil {
		panic(err)
	}
}

func clearReportFile(){
	d1 := []byte("")
	err := ioutil.WriteFile("./report", d1, 0644)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
