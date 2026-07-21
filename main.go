package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
)

const (
	broadcast = "255.255.255.255"
	mac = "112233445566"
)

func main() {
	conn, err := net.Dial("udp", broadcast+":9")
	if err != nil {
		fmt.Println("Error sending packet", err)
		return
	}
	defer conn.Close()

	synchronization := bytes.Repeat([]byte{0xFF}, 6)

	macBytes, _ := hex.DecodeString(mac)

	repeatedMac := bytes.Repeat(macBytes, 16)

	conn.Write(append(synchronization, repeatedMac...))
}
