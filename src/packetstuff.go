package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/shirou/gopsutil/net"
	"os/exec"
	"strings"
)

type appPackets struct {
	Name       string
	IP         string
	Port       int
	packetList []gopacket.Packet
}

func showpackets(packetSource *gopacket.PacketSource) {
	appsPakketList := []appPackets{}
	go func() {
		for packet := range packetSource.Packets() {
			if !isPaused {
				fmt.Println("----------------------------------------")

				for _, layer := range packet.Layers() {
					if layer.LayerType() != layers.LayerTypeIPv4 && layer.LayerType() != layers.LayerTypeIPv6 && layer.LayerType() != layers.LayerTypeTCP && layer.LayerType() != layers.LayerTypeUDP && layer.LayerType() != layers.LayerTypeEthernet && layer.LayerType().String() != "Payload" {
						fmt.Println("PACKET LAYER:", layer.LayerType())
					}
				}

				ipv4Layer := packet.Layer(layers.LayerTypeIPv4)

				otherIp := ""
				appPort := 0
				if ipv4Layer != nil {
					ipv4, _ := ipv4Layer.(*layers.IPv4)
					if ipv4 != nil {
						tcpLayer := packet.Layer(layers.LayerTypeTCP)
						udpLayer := packet.Layer(layers.LayerTypeUDP)
						if ipv4.SrcIP.String() == ipv4AddrOfInterface {
							otherIp = ipv4.DstIP.String()
							if tcpLayer != nil {
								tmptcp, _ := tcpLayer.(*layers.TCP)
								appPort = int(tmptcp.DstPort)
							}
							if udpLayer != nil {
								udpLayer, _ := udpLayer.(*layers.UDP)
								appPort = int(udpLayer.DstPort)
							}
						} else if ipv4.DstIP.String() == ipv4AddrOfInterface {
							otherIp = ipv4.SrcIP.String()
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

				appName := getAppName(appPort)
				fmt.Println("myIp: ", ipv4AddrOfInterface)
				fmt.Println("otherIp: ", otherIp)
				fmt.Println("appPort: ", appPort)
				fmt.Println("appName ", appName)

				isNotInList := true
				for _, appPakket := range appsPakketList {
					if appPakket.Name == appName && appPakket.IP == otherIp && appPakket.Port == appPort {
						isNotInList = false
						break
					}
				}

				if isNotInList {
					appsPakketList = append(appsPakketList, appPackets{
						Name:       appName,
						IP:         otherIp,
						Port:       appPort,
						packetList: []gopacket.Packet{packet},
					})
				} else {
					for _, appPakket := range appsPakketList {
						if appPakket.Name == appName && appPakket.IP == otherIp && appPakket.Port == appPort {
							appPakket.packetList = append(appPakket.packetList, packet)
							break
						}
					}
				}

				fmt.Println("appsPakketList: ", len(appsPakketList))
			}
		}
	}()

}

func getAppName(port int) string {
	//fmt.Println("port: ", port)

	connections, err := net.Connections("all")
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	pid := int32(0)
	for _, conn := range connections {
		if int(conn.Laddr.Port) == port || int(conn.Raddr.Port) == port {
			pid = conn.Pid
			break
		}
	}
	return getProcessRunningStatus(int(pid))
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
