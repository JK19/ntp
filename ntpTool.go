package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/JK19/ntp/ntpPacket"
)

func main() {

	serverFlag := flag.String("s", "", "NTP server address")

	zoneFlag := flag.String("tz", "Local", "Timezone from IANA timezone database")

	leapFlag := flag.Uint("li", 0, "Leap indicator")

	verFlag := flag.Uint("ver", 3, "NTP protocol version")

	flag.Usage = func() {
		fmt.Println("\nUsage: ntp [-s <server>] [-tz <timezone>] [-ver <version>] [-li <indicator>] ")
		fmt.Println("\n\nOptions: ")
		fmt.Println("-s <server>		NTP server address")
		fmt.Println("-tz <timezone>		Timezone from IANA timezone database")
		fmt.Println("-ver <version>		NTP protocol version")
		fmt.Println("-li <indicator>		0 -No leap second adjustment")
		fmt.Println("			1 -Last minute of the day has 61 seconds")
		fmt.Println("			2 -Last minute of the day has 59 seconds")
		fmt.Println("			3 -Clock is unsynchronized")
		fmt.Println("\nDefaults: ")
		fmt.Println("-tz Local")
		fmt.Println("-ver 3")
		fmt.Println("-li 0")

	}

	flag.Parse()

	if *serverFlag == "" {
		fmt.Fprint(os.Stderr, "ERROR: a server must be specified\n")
		flag.Usage()
		os.Exit(1)
	}

	loc, err := time.LoadLocation(*zoneFlag)

	if err != nil {
		fmt.Fprint(os.Stderr, "ERROR: could not find the timezone specified\n")
		os.Exit(1)
	}

	address := *serverFlag + ":" + "123"

	fmt.Println("\nRequesting time to: " + address + "\n")

	socket, err := net.Dial("udp", address)

	if err != nil {
		fmt.Fprint(os.Stderr, "ERROR: could not create socket\n")
		os.Exit(1)
	}

	defer socket.Close()

	packet := ntpPacket.NewNtpPacket()

	packet.SetLeap(uint8(*leapFlag))
	packet.SetVersion(uint8(*verFlag))

	err = packet.SendTo(socket)

	if err != nil {
		fmt.Fprint(os.Stderr, "ERROR: failed to send data to server\n")
		os.Exit(1)
	}

	err = packet.ReadFrom(socket)

	if err != nil {
		fmt.Fprint(os.Stderr, "ERROR: failed to read data from server\n")
		os.Exit(1)
	}

	fmt.Printf("--- Response from %s ---\n", *serverFlag)

	fmt.Println("Leap indicator: ", packet.GetLeap())
	fmt.Println("Version: ", packet.GetVersion())
	fmt.Println("Mode: ", packet.GetMode())
	fmt.Println("Stratum: ", packet.Getstratum())
	fmt.Println("Poll interval: ", packet.GetPollInterval())
	fmt.Println("Precision: ", packet.Getprecision())
	fmt.Println("Root delay: ", packet.GetRootDelay())
	fmt.Println("Root dispersion: ", packet.GetRootDispersion())
	fmt.Println("Reference clock id: ", packet.GetRefClokId())
	fmt.Println("Reference timestamp: ", packet.GetRefTimestamp())
	fmt.Println("Originate timestamp: ", packet.GetOriginTimestamp())
	fmt.Println("Receive timestamp: ", packet.GetRxTimestamp())
	fmt.Println("Transmit timestamp: ", packet.GetTxTimestamp())

	fmt.Println("\n\nTime from server: ", packet.GetTime().In(loc))
}
