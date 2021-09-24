package main

import "fmt"

type IPAddr [4]byte

// TODO: Add a "String() string" method to IPAddr.
func (ip IPAddr) String() string{
	var newIPAddress string
	for _, num := range ip {
		if newIPAddress!="" {
			newIPAddress += "."
		}
		newIPAddress+=fmt.Sprintf("%v",num)
	}	
	return newIPAddress
}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
}
