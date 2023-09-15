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
	counter   gopacket.Packet
	counterMu sync.Mutex
	clients   []chan gopacket.Packet
	clientMu  sync.Mutex
)

func main() {

	// Find all available network devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	// Print information about each network device
	for _, device := range devices {
		// Filter out virtual network devices
		if !(strings.Contains(device.Description, "VMnet")) && !(strings.Contains(device.Description, "Virtual")) {
			log.Printf("Name: %s\nDescription: %s\n", device.Name, device.Description)
		}
	}

	// Open a network interface for capturing packets
	handle, err := pcap.OpenLive("\\Device\\NPF_Loopback", 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	err = handle.SetBPFFilter("tcp") // Set a BPF filter to capture only TCP packets
	if err != nil {
		log.Fatal(err)
	}

	// Create a packet source from the handle
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./src/index.html") // Specify the correct path to your HTML file
	})

	http.HandleFunc("/updatePackets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Create a channel for this client
		clientChan := make(chan gopacket.Packet)

		// Register the client
		clientMu.Lock()
		clients = append(clients, clientChan)
		clientMu.Unlock()

		// Notify when the client's connection is closed
		closeNotifier := w.(http.CloseNotifier).CloseNotify()

		go func() {
			<-closeNotifier
			// Handle client disconnect here
			removeClient(clientChan)
		}()

		for {
			// Send updates to the client
			select {
			case counterValue := <-clientChan:
				fmt.Fprintf(w, "data: %d\n\n", counterValue)
				w.(http.Flusher).Flush()
			case <-r.Context().Done():
				// Handle client disconnect here
				removeClient(clientChan)
				return
			}
		}
	})
	go func() {

		// Process incoming packets
		for packet := range packetSource.Packets() {
			// Process the outgoing packet here
			fmt.Println(packet)
			updateClients(packet)
		}
	}()
	http.ListenAndServe(":8888", nil)

	/*	for {
			time.Sleep(1 * time.Second)
			incrementCounter()
			updateClients(counter)
		}
	*/

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
			// Remove the client from the list
			clients = append(clients[:i], clients[i+1:]...)
			close(clientChan) // Close the client's channel
			break
		}
	}
}
