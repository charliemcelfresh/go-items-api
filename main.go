// Package main instantiates and runs our simple JSON api server
package main

import "github.com/charliemcelfresh/go-items-api/internal/server"

func main() {
	s := server.NewServer()
	s.Run()
}
