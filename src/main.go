package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

var (
	counter     gopacket.Packet
	counterMu   sync.Mutex
	clients     []chan gopacket.Packet
	clientMu    sync.Mutex
	isPaused    bool
	pauseResume sync.Mutex
	handle      *pcap.Handle
)

func main() {

	http.HandleFunc("/", handleSummonWebpage)
	http.HandleFunc("/updatePackets", handleUpdatePackets)
	http.HandleFunc("/pause", handlePause)
	http.HandleFunc("/resume", handleResume)
	http.HandleFunc("/selectDevice", handleSelectDevice)

	http.ListenAndServe(":8888", nil)

}
func getHostIPAddress() string {
	// Get the host's IP address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	log.Fatal("Unable to determine host's IP address")
	return ""
}
func handleSelectDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Println(handle)
	if handle != nil {
		handle.Close()
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	device := string("")
	for _, value := range body {
		asciiChar := fmt.Sprintf("%c", value)

		device += asciiChar
	}
	//test := "\\Device\\NPF_Loopback"
	//fmt.Println(test)

	fmt.Println(device)
	handle, err := pcap.OpenLive(device, 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	err = handle.SetBPFFilter("tcp and port 80") // Capture only HTTP traffic (port 80)
	if err != nil {
		log.Fatal(err)
	}
	//hostIP := getHostIPAddress()
	//hostIP := "192.168.1.105"
	//fmt.Println(hostIP)
	//err = handle.SetBPFFilter("tcp and (src host " + hostIP + " or dst host " + hostIP + ")")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = handle.SetBPFFilter("tcp")
	//err = handle.SetBPFFilter("udp")
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Send a response to the client
	fmt.Fprintln(w, "Device selected: "+device)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	go func() {
		for packet := range packetSource.Packets() {
			if !isPaused {
				//fmt.Println(packet)
				updateClients(packet)
			}
		}
	}()

	//defer handle.Close()
}

func updateClients(counter gopacket.Packet) {
	clientMu.Lock()
	defer clientMu.Unlock()
	for _, clientChan := range clients {
		clientChan <- counter
	}
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

func handleUpdatePackets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientChan := make(chan gopacket.Packet)

	clientMu.Lock()
	clients = append(clients, clientChan)
	clientMu.Unlock()

	closeNotifier := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-closeNotifier
		removeClient(clientChan)
	}()

	for {
		select {
		case counterValue := <-clientChan:
			fmt.Fprintf(w, "data: %d\n\n", counterValue)
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			removeClient(clientChan)
			return
		}
	}
}

func handleSummonWebpage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./src/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Buttons []string
	}{
		Buttons: selectAbleDevices(),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func selectAbleDevices() []string {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	var deviceNames []string

	for _, device := range devices {
		if !(strings.Contains(device.Description, "VMnet")) && !(strings.Contains(device.Description, "Virtual")) && !(strings.Contains(device.Description, "Bluetooth")) && !(strings.Contains(device.Description, "Miniport")) {
			fmt.Println(device.Description)
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
