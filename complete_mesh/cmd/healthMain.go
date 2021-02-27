package main

import (
	healthcheck "github.com/Open-Twin/citymesh/complete_mesh/healtchcheck"
)

func main() {

	//healthcheck.Healthcheck()
	ids := []string{"Master123", "Master456"}
	ips := []string{"192.168.0.2", "192.168.0.1"}
	healthcheck.Healthcheck(ids, ips)
}