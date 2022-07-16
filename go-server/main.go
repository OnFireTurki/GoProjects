package main

import "time"

func main() {
	// Run the server and send requests using client

	go func() {
		// Giving the server time to run and it must run in the main thread so the program does not close
		// Using wait groups in this case is pointless, it is just simple app
		time.Sleep(time.Second * 3)
		client()
	}()
	server()
}
