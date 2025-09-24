package main

import (
	"HttpToTcp/internal/request"
	"fmt"
	"io"
	"net"
	"os"
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

		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Request Line:")
		fmt.Println(" - Method: ", r.RequestLine.Method)
		fmt.Println(" - Target: ", r.RequestLine.RequestTarget)
		fmt.Println(" - Version: ", r.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range r.Headers {
			fmt.Printf(" - %s: %s\n", k, v)
		}
		fmt.Println("Body: ")
		fmt.Println(" - Body_String: ", string(r.Body))
	}

}
