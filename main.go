package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	content := make(chan string)

	reader := make([]byte, 8)
	holder := ""

	go func() {
		for {
			c, err := f.Read(reader)
			if err != nil {
				break
			}

			holder += string(reader[:c])

		}
		content <- holder
	}()
	return content
}

func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Connection Accepted")

		context := getLinesChannel(conn)
		val, ok := <-context
		if !ok {
			fmt.Println("Connection Closed")
		}

		lines := strings.SplitSeq(val, "\n")
		for l := range lines {
			fmt.Printf("read: %s\n", l)
		}
	}

}
