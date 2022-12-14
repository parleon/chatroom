package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const SERVER = "xS3Ver_@m1N"

func initialize_outgoing(host string, port string) net.Conn {
	conn, err := net.Dial("tcp", host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func main() {
	host := os.Args[1]
	port := os.Args[2]
	username := os.Args[3]

	conn := initialize_outgoing(host, port)

	enc := gob.NewEncoder(conn)
	
	enc.Encode(map[string]string{
		"to":SERVER,
		"from": username,
		"message": "faueuyfhiwdufh",
	})

	go func() {
		dec := gob.NewDecoder(conn)
		message_map := make(map[string]string)
		for{
			decerr := dec.Decode(&message_map)
			if decerr != nil {
				fmt.Println("decerr")
				return
			}
			fmt.Println(message_map)
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		fields := strings.Fields(text)

		message := strings.Join(fields[1:], " ")

		go func ()  {
			outgoing := make(map[string]string)
			outgoing["from"] = username
			outgoing["to"] = fields[0]
			outgoing["message"] = message
			enc.Encode(outgoing)
		}()

	}
	
}