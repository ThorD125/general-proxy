package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type appPackets struct {
	Name       string
	IP         string
	Port       int
	packetList []gopacket.Packet
}

func showpackets(packetSource *gopacket.PacketSource) {
	//var appsPakketList = []appPackets{}
	//appsPakketList := make(map[string][]gopacket.Packet)
	appsPakketList := make(map[string][][]byte)
	go func() {
		for packet := range packetSource.Packets() {
			if !isPaused {

				for _, layer := range packet.Layers() {
					if layer.LayerType() != layers.LayerTypeIPv4 && layer.LayerType() != layers.LayerTypeIPv6 && layer.LayerType() != layers.LayerTypeTCP && layer.LayerType() != layers.LayerTypeUDP && layer.LayerType() != layers.LayerTypeEthernet && layer.LayerType().String() != "Payload" {
						fmt.Println("PACKET LAYER:", layer.LayerType())
					}
				}

				ipv4Layer := packet.Layer(layers.LayerTypeIPv4)

				appPort := 0
				if ipv4Layer != nil {
					ipv4, _ := ipv4Layer.(*layers.IPv4)
					if ipv4 != nil {
						tcpLayer := packet.Layer(layers.LayerTypeTCP)
						udpLayer := packet.Layer(layers.LayerTypeUDP)
						if ipv4.SrcIP.String() == ipv4AddrOfInterface {
							if tcpLayer != nil {
								tmptcp, _ := tcpLayer.(*layers.TCP)
								appPort = int(tmptcp.DstPort)
							}
							if udpLayer != nil {
								udpLayer, _ := udpLayer.(*layers.UDP)
								appPort = int(udpLayer.DstPort)
							}
						} else if ipv4.DstIP.String() == ipv4AddrOfInterface {
							if tcpLayer != nil {
								tmptcp, _ := tcpLayer.(*layers.TCP)
								appPort = int(tmptcp.SrcPort)
							}
							if udpLayer != nil {
								udpLayer, _ := udpLayer.(*layers.UDP)
								appPort = int(udpLayer.SrcPort)
							}
						}
					}
				}

				appName := getProcessRunningStatus(getAppName(appPort))
				appsPakketList[appName] = append(appsPakketList[appName], packet.Data())

				fmt.Println(packet)
				//updateClients(appsPakketList)

				fmt.Println("----------------------------------------")
				fmt.Println("somanypackets: ", len(appsPakketList))
				for appName, Pakket := range appsPakketList {
					fmt.Println(appName, ": ", len(Pakket))
				}

			} else {
				break
			}
		}
	}()

}
