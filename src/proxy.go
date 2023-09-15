// Import necessary packages
package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"strings"
)

// Define a struct to represent a network device
type device struct {
	name        string
	description string
}

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

	// Process incoming packets
	for packet := range packetSource.Packets() {
		// Process the outgoing packet here
		fmt.Println(packet)
	}

	defer handle.Close()
}
