package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
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
)

func main() {
	handle, err := pcap.OpenLive(selectDevice(), 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	err = handle.SetBPFFilter("tcp")
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	http.HandleFunc("/", handleSummonWebpage)
	http.HandleFunc("/updatePackets", handleUpdatePackets)
	http.HandleFunc("/pause", handlePause)
	http.HandleFunc("/resume", handleResume)

	go func() {
		for packet := range packetSource.Packets() {
			if !isPaused {
				//fmt.Println(packet)
				updateClients(packet)
			}
		}
	}()
	http.ListenAndServe(":8888", nil)

	defer handle.Close()
}

func getCounter() gopacket.Packet {
	counterMu.Lock()
	defer counterMu.Unlock()
	return counter
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
	http.ServeFile(w, r, "./src/index.html")
}

func selectDevice() string {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		if !(strings.Contains(device.Description, "VMnet")) && !(strings.Contains(device.Description, "Virtual")) {
			log.Printf("Name: %s\nDescription: %s\n", device.Name, device.Description)
		}
	}
	return "\\Device\\NPF_Loopback"
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
