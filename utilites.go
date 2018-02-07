package main

import (
	"time"
	"fmt"
	"math/rand"
	"io/ioutil"
)

func getRandomInt(max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(max)
}

func printHex(data []byte){
	fmt.Printf("[ ")
	for _, i := range data {
		fmt.Printf("0x%x, ", i)
	}
	fmt.Println("]")

}

func logToConsole(message string){
	fmt.Println(time.Now(), message)
}

func logToFile(message string){
	filename := "/Users/jon/log.txt"
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	writeMessage := fmt.Sprintf("%s\n%s", string(d), message)
	fmt.Println(writeMessage)
	ioutil.WriteFile(filename, []byte(writeMessage), 0666)
}