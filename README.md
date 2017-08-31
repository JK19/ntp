# NTP Packet

**NTP Packet** is a library for querying NTP servers that allows access to all fields of the ntp packet and read them easily. It comes with **ntp**, a command line tool that allows you to perform NTP queries from the terminal.

## Installation
In order o install **Ntp Packet** you need Go 1.8 or later, and a configured Go enviroment.
To install, run:

    go install github.com/JK19/ntp

That will install the library and the **ntp** tool

To install only the library, run:

    go install github.com/JK19/ntp/ntpPacket

## Usage

Creating a NTP packet and quering the server:
```
package main

import (
	"fmt"
	"net"
	"time"

	"github.com/JK19/ntp/ntpPacket"
)

func main() {

	server := "3.es.pool.ntp.org:123"

	socket, _ := net.Dial("udp", server)

	defer socket.Close()

	packet := ntpPacket.NewNtpPacket()

	packet.SendTo(socket)

	packet.ReadFrom(socket)

	fmt.Println("NTP time: ", packet.GetTime())
	fmt.Println("Local time: ", packet.GetTime().In(time.Local))
}
```
For more details about usage check out the **ntp** tool implementation in [ntpTool.go](ntpTool.go)


## ntp tool

> in order to use **ntp** from terminal you should have your $GOPATH/bin added to your $PATH env variable (or %GOPATH%\bin in windows)

To query a ntp server from terminal, run:

    ntp -s 3.es.pool.ntp.org

It will give you the fields from the ntp packet and the time (local by default):

```
Requesting time to: 3.es.pool.ntp.org:123

--- Response from 3.es.pool.ntp.org ---
Leap indicator:  0
Version:  3
Mode:  4
Stratum:  1
Poll interval:  8s
Precision:  119ns
Root delay:  0s
Root dispersion:  6s
Reference clock id:  GPS
Reference timestamp:  2017-08-31 16:16:00.189979106 +0000 UTC
Originate timestamp:  1900-01-01 00:00:00 +0000 UTC
Receive timestamp:  2017-08-31 16:16:11.698723434 +0000 UTC
Transmit timestamp:  2017-08-31 16:16:11.698765314 +0000 UTC


Time from server:  2017-08-31 18:16:11.698765314 +0200 CEST
```
You can set the parameters of the query, to get the available parameters run:

    ntp -h


```
Usage: ntp [-s <server>] [-tz <timezone>] [-ver <version>] [-li <indicator>]


Options:
-s <server>             NTP server address
-tz <timezone>          Timezone from IANA timezone database
-ver <version>          NTP protocol version
-li <indicator>         0 -No leap second adjustment
                        1 -Last minute of the day has 61 seconds
                        2 -Last minute of the day has 59 seconds
                        3 -Clock is unsynchronized

Defaults:
-tz Local
-ver 3
-li 0
```

The timezones follow the [IANA timezone database](https://www.iana.org/time-zones) format [[wikipedia]](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).
Timezone examples:

* Europe/Madrid
* America/New_York
* Europe/Moscow

    ntp -s 3.es.pool.ntp.org -tz Europe/Madrid

## License
* [MIT](LICENSE)

