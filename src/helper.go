package main

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"github.com/shirou/gopsutil/net"
	"log"
	"net/http"
	"os/exec"
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

func removeClient(clientChan chan map[string][][]byte) {
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

func getAppName(port int) int {

	connections, err := net.Connections("all")
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	pid := int32(0)
	for _, conn := range connections {
		if int(conn.Laddr.Port) == port || int(conn.Raddr.Port) == port {
			pid = conn.Pid
			break
		}
	}
	return int(pid)
}

func getProcessRunningStatus(pid int) string {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV")

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	exeName := ""
	if len(lines) >= 2 {
		exeName = strings.Trim(strings.Split(lines[1], ",")[0], "\"")
	}
	return exeName
}
func updateClients(packetsing map[string][][]byte) {
	clientMu.Lock()
	defer clientMu.Unlock()
	for _, clientChan := range clients {
		clientChan <- packetsing
	}
}
