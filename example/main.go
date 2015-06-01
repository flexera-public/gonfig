package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Couldn't retrieve config path")
	}
	exampleDir := filepath.Dir(filename)
	err := Load(filepath.Join(exampleDir, "config.json"))
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}
	if ApiEndpoint != "https://us-3.rightscale.com" {
		panic("Expected us-3.rightscale.com for ApiEndpoint value but got " + ApiEndpoint)
	}
	if Port != int64(8000) {
		panic("Expected 8000 for Port value but got " + strconv.Itoa(int(Port)))
	}
	if Worker == nil {
		panic("Exepected non-nil value for Worker but got nil")
	}
	if Worker.Concurrency != 20 {
		panic("Exepected 20 for Worker.Concurrency but got " + strconv.Itoa(int(Worker.Concurrency)))
	}
	if !Worker.Enabled {
		panic("Exepected Worker.Enabled to be true but got false")
	}
	fmt.Println("ok")
}
