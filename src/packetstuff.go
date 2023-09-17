package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type packetDst struct {
	app        string
	appPort    int
	packetList []gopacket.Packet
}

func showpackets(packetSource *gopacket.PacketSource) {
	go func() {
		for packet := range packetSource.Packets() {
			if !isPaused {
				fmt.Printf("----------------------------------------\n")

				fmt.Println("ip: ", ipAddrOfInterface)

				for _, layer := range packet.Layers() {
					fmt.Println("PACKET LAYER:", layer.LayerType())

					ipLayer := packet.Layer(layers.LayerTypeIPv4)
					ipv6Layer := packet.Layer(layers.LayerTypeIPv6)
					tcpLayer := packet.Layer(layers.LayerTypeTCP)
					udpLayer := packet.Layer(layers.LayerTypeUDP)

					if ipLayer != nil {
						ip, _ := ipLayer.(*layers.IPv4)
						fmt.Println("IP SrcIP:", ip.SrcIP)
						fmt.Println("IP DstIP:", ip.DstIP)
					}
					if ipv6Layer != nil {
						ipv6, _ := ipv6Layer.(*layers.IPv6)
						fmt.Println("IPv6 SrcIP:", ipv6.SrcIP)
						fmt.Println("IPv6 DstIP:", ipv6.DstIP)
					}
					if tcpLayer != nil {
						tmptcp, _ := tcpLayer.(*layers.TCP)
						fmt.Println("TCP SrcPort:", tmptcp.SrcPort)
						fmt.Println("TCP DstPort:", tmptcp.DstPort)
					}
					if udpLayer != nil {
						udpLayer, _ := udpLayer.(*layers.UDP)
						fmt.Println("UDP SrcPort:", udpLayer.SrcPort)
						fmt.Println("UDP DstPort:", udpLayer.DstPort)
					}

					//ethLayer := packet.Layer(layers.LayerTypeEthernet)
					//if ethLayer != nil {
					//	ethPacket, _ := ethLayer.(*layers.Ethernet)
					//	fmt.Println("Ethernet Source MAC:", ethPacket.SrcMAC)
					//	fmt.Println("Ethernet Destination MAC:", ethPacket.DstMAC)
					//	fmt.Println("Ethernet Ethertype:", ethPacket.EthernetType)
					//	//fmt.Println("Ethernet Contents:", ethPacket.Contents)
					//	//fmt.Println("Ethernet Payload:", ethPacket.Payload)
					//}
					//if ipLayer != nil {
					//	ip, _ := ipLayer.(*layers.IPv4)
					//	//---
					//	// flag is 3 bits long
					//	// first bit is always 0
					//	// second bit is DF (Don't Fragment) bit
					//	// third bit is MF (More Fragments) bit
					//	fmt.Println("IP Flags:", ip.Flags)
					//	//FragOffset is 13 bits long
					//	//FragOffset is the offset of the data in the original datagram, measured in units of 8 octets (64 bits).
					//	fmt.Println("IP FragOffset:", ip.FragOffset)
					//	//---
					//	fmt.Println("IP Version:", ip.Version)
					//	fmt.Println("IP Protocol:", ip.Protocol)
					//	fmt.Println("IP SrcIP:", ip.SrcIP)
					//	fmt.Println("IP DstIP:", ip.DstIP)
					//	//fmt.Println("IP Contents:", ip.Contents)
					//	//fmt.Println("IP Payload:", ip.Payload)
					//}
					//if ipv6Layer != nil {
					//	ipv6, _ := ipv6Layer.(*layers.IPv6)
					//	fmt.Println("IPv6 Version:", ipv6.Version)
					//	fmt.Println("IPv6 NextHeader:", ipv6.NextHeader)
					//	fmt.Println("IPv6 SrcIP:", ipv6.SrcIP)
					//	fmt.Println("IPv6 DstIP:", ipv6.DstIP)
					//	//fmt.Println("IPv6 Payload:", ipv6.Payload)
					//}
					//if tcpLayer != nil {
					//	tmptcp, _ := tcpLayer.(*layers.TCP)
					//	fmt.Println("TCP SrcPort:", tmptcp.SrcPort)
					//	fmt.Println("TCP DstPort:", tmptcp.DstPort)
					//}
					if udpLayer != nil {
						udpLayer, _ := udpLayer.(*layers.UDP)
						fmt.Println("UDP SrcPort:", udpLayer.SrcPort)
						fmt.Println("UDP DstPort:", udpLayer.DstPort)
					}
				}
			}
		}
	}()

}
