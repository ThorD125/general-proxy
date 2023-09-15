package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	counter   int
	counterMu sync.Mutex
	clients   []chan int
	clientMu  sync.Mutex
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./src/index.html") // Specify the correct path to your HTML file
	})

	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Create a channel for this client
		clientChan := make(chan int)

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
		for {
			time.Sleep(1 * time.Second)
			incrementCounter()
			updateClients(counter)
		}
	}()

	http.ListenAndServe(":8888", nil)
}

func incrementCounter() {
	counterMu.Lock()
	defer counterMu.Unlock()
	counter++
}

func getCounter() int {
	counterMu.Lock()
	defer counterMu.Unlock()
	return counter
}

func updateClients(counter int) {
	clientMu.Lock()
	defer clientMu.Unlock()
	for _, clientChan := range clients {
		clientChan <- counter
	}
}

func removeClient(clientChan chan int) {
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
