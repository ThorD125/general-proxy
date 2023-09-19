package main

import (
	"github.com/google/gopacket/pcap"
	"sync"
)

var (
	clients     []chan map[string][][]byte
	clientMu    sync.Mutex
	pauseResume sync.Mutex
	handle      *pcap.Handle

	ipv4AddrOfInterface string
)

var isPaused = true
