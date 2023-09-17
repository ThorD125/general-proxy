package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"net/http"
	"strings"
)

func selectAbleDevices() []string {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	var deviceNames []string

	for _, device := range devices {
		if !(strings.Contains(device.Description, "VMnet")) && !(strings.Contains(device.Description, "Virtual")) && !(strings.Contains(device.Description, "Bluetooth")) && !(strings.Contains(device.Description, "Miniport")) {
			//fmt.Println(device.Description)
			//fmt.Println(device.Addresses[0].IP.String())
			deviceNames = append(deviceNames, device.Name)
		}
	}

	return deviceNames
}

func handlePause(w http.ResponseWriter, r *http.Request) {
	pauseResume.Lock()
	isPaused = true
	pauseResume.Unlock()
	fmt.Fprintln(w, "Capture Paused")
}

func handleResume(w http.ResponseWriter, r *http.Request) {
	pauseResume.Lock()
	isPaused = false
	pauseResume.Unlock()
	fmt.Fprintln(w, "Capture Resumed")
}

func numbersToASCII(numbers []int) string {
	var asciiString string

	for _, num := range numbers {
		// Convert the number to its ASCII representation and append it to the string
		asciiChar := string(num)
		asciiString += asciiChar
	}

	return asciiString
}

func removeClient(clientChan chan gopacket.Packet) {
	clientMu.Lock()
	defer clientMu.Unlock()
	for i, c := range clients {
		if c == clientChan {
			clients = append(clients[:i], clients[i+1:]...)
			close(clientChan)
			break
		}
	}
}
