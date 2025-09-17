package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	dest, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, dest)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		conn.Write([]byte(line))
	}
}
