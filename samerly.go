package main

import (
	"github.com/samertm/samerly/server"
	
	"fmt"
)

var _ = fmt.Println // debugging

func main() {
	ip := "localhost:8520"
	fmt.Println("Listening on", ip)
	server.ListenAndServe(ip)
}
