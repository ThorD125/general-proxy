package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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

	ethLayer layers.Ethernet
	ipLayer  layers.IPv4
	tcpLayer layers.TCP
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

type Payload struct {
	Field1 uint8
	Field2 uint8
	Field3 uint8
	Field4 uint8
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

	if body == nil {
		device = "\\Device\\NPF_Loopback"
	}
	//test := "\\Device\\NPF_Loopback"
	//fmt.Println(test)

	fmt.Println(device)

	handle, err := pcap.OpenLive(device, 65536, true, pcap.BlockForever)

	if err != nil {
		log.Fatal(err)
	}
	//err = handle.SetBPFFilter("tcp and port 80") // Capture only HTTP traffic (port 80)
	//if err != nil {
	//	log.Fatal(err)
	//}
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
		/*for packet := range packetSource.Packets() {
			if !isPaused {

				fmt.Println(packet)
				updateClients(packet)
			}
		}*/
		for packet := range packetSource.Packets() {
			if !isPaused {
				fmt.Printf("----------------------------------------\n")

				for _, layer := range packet.Layers() {
					fmt.Println("PACKET LAYER:", layer.LayerType())
					ethLayer := packet.Layer(layers.LayerTypeEthernet)
					if ethLayer != nil {
						//ethPacket, _ := ethLayer.(*layers.Ethernet)
						//fmt.Println("Ethernet Source MAC:", ethPacket.SrcMAC)
						//fmt.Println("Ethernet Destination MAC:", ethPacket.DstMAC)
						//fmt.Println("Ethernet Ethertype:", ethPacket.EthernetType)
						//fmt.Println("Ethernet Contents:", ethPacket.Contents)
						//fmt.Println("Ethernet Payload:", ethPacket.Payload)
					}

					ipLayer := packet.Layer(layers.LayerTypeIPv4)
					if ipLayer != nil {
						//ip, _ := ipLayer.(*layers.IPv4)
						//---
						// flag is 3 bits long
						// first bit is always 0
						// second bit is DF (Don't Fragment) bit
						// third bit is MF (More Fragments) bit
						//fmt.Println("IP Flags:", ip.Flags)
						// FragOffset is 13 bits long
						// FragOffset is the offset of the data in the original datagram, measured in units of 8 octets (64 bits).
						//fmt.Println("IP FragOffset:", ip.FragOffset)
						//---
						//fmt.Println("IP Version:", ip.Version)
						//fmt.Println("IP Protocol:", ip.Protocol)
						//fmt.Println("IP SrcIP:", ip.SrcIP)
						//fmt.Println("IP DstIP:", ip.DstIP)
						//fmt.Println("IP Contents:", ip.Contents)
						//fmt.Println("IP Payload:", ip.Payload)
					}
					ipv6Layer := packet.Layer(layers.LayerTypeIPv6)
					if ipv6Layer != nil {
						//ipv6, _ := ipv6Layer.(*layers.IPv6)

						//fmt.Println("IPv6 Version:", ipv6.Version)
						//fmt.Println("IPv6 NextHeader:", ipv6.NextHeader)
						//fmt.Println("IPv6 SrcIP:", ipv6.SrcIP)
						//fmt.Println("IPv6 DstIP:", ipv6.DstIP)
						//fmt.Println("IPv6 Payload:", ipv6.Payload)
					}
					tcpLayer := packet.Layer(layers.LayerTypeTCP)
					if tcpLayer != nil {
						//tmptcp, _ := tcpLayer.(*layers.TCP)
						//fmt.Println("TCP SrcPort:", tmptcp.SrcPort)
						//might be used to check what app is using it
						//fmt.Println("TCP DstPort:", tmptcp.DstPort)
					}
					udpLayer := packet.Layer(layers.LayerTypeUDP)
					if udpLayer != nil {
						//udpLayer, _ := udpLayer.(*layers.UDP)
						//fmt.Println("UDP SrcPort:", udpLayer.SrcPort)
						//fmt.Println("UDP DstPort:", udpLayer.DstPort)
						//fmt.Println("UDP Payload:", udpLayer.Payload)
					}
					payloadLayer := packet.ApplicationLayer()
					if payloadLayer != nil {
						payloadLayer, _ := payloadLayer.(*gopacket.Payload)
						fmt.Println("Payload:", payloadLayer.Payload())

						var payload Payload
						fmt.Println("Payload:", binary.Read(bytes.NewReader(payloadLayer.Payload()), binary.LittleEndian, &payload))

					}
				}
			}
		}

	}()

	//defer handle.Close()
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
