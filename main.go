package main

import (
	"flag"
	"log"
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

func main() {
	var (
		ifaceFlag = flag.String("i", "", "network interface to use to send and receive messages")
	)

	flag.Parse()

	log.Printf("Shinfuku1 Start")

	ifi, err := net.InterfaceByName(*ifaceFlag)

	if err != nil {
		log.Fatalf("failed to open interface %q: %v", *ifaceFlag, err)
	}

	// Open a raw socket using same EtherType as our frame.
	c, err := raw.ListenPacket(ifi, 0x0800, nil)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer c.Close()

	// Accept frames up to interface's MTU in size.
	b := make([]byte, ifi.MTU)
	var f ethernet.Frame

	log.Printf("Start capture. ")

	// Keep reading frames.
	for {
		n, addr, err := c.ReadFrom(b)
		if err != nil {
			log.Fatalf("failed to receive message: %v", err)
		}

		// Unpack Ethernet frame into Go representation.
		if err := (&f).UnmarshalBinary(b[:n]); err != nil {
			log.Fatalf("failed to unmarshal ethernet frame: %v", err)
		}

		// Display source of message and message itself.
		log.Printf("[%s] %X", addr.String(), string(f.Payload))
	}

}
