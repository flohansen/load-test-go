package main

import (
	"flag"
	"load-test/load"
	"load-test/net"
	"log"
)

func main() {
	numberUsers := flag.Int("users", 500, "The number of users/threads")
	rampupTimeSeconds := flag.Int("rampup", 60, "The rampup time until the full load takes effect (in seconds)")
	durationSeconds := flag.Int("duration", 300, "The duration of the load test (in seconds)")
	target := flag.String("target", "localhost:80", "The target to test in format host:port")
	path := flag.String("path", "/", "The path to test")
	flag.Parse()

	config := load.TesterConfig{
		NumberUsers: *numberUsers,
		Rampup:      *rampupTimeSeconds,
		Duration:    *durationSeconds,
		Target:      *target,
		Path:        *path,
	}

	client := net.NewSocketClient()
	tester := load.NewTester(config, client)
	if err := tester.Run(); err != nil {
		log.Fatal(err)
	}
}
