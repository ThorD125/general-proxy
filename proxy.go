package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	// Print information about each network device
	for _, device := range devices {
		log.Printf("Name: %s\nDescription: %s\n", device.Name, device.Description)
	}

	handle, err := pcap.OpenLive("\\Device\\NPF_Loopback", 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	err = handle.SetBPFFilter("tcp")
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Process the outgoing packet here
		fmt.Println(packet)
	}

	defer handle.Close()
}
