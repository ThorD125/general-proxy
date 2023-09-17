package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

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

	showpackets(packetSource)

	//defer handle.Close()
}
