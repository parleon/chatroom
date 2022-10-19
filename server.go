package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
)

const SERVER = "xS3Ver_@m1N"

func initialize_source(port string) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatal(err)
	}
	return ln
}

type message_wrapper struct {
	message     map[string]string
	sender_conn net.Conn
}

func main() {
	port := os.Args[1]

	source := initialize_source(port)
	m_queue := make(chan message_wrapper, 500)

	go func() {
		for {
			conn, err := source.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go func() {
				dec := gob.NewDecoder(conn)
				message_map := make(map[string]string)
				for {
					dec.Decode(&message_map)
					m_queue <- message_wrapper{message: message_map, sender_conn: conn}
				}
			}()
		}
	}()

	user_connections := make(map[string]net.Conn)
	for {
		select {
		case m := <-m_queue:
			fmt.Println(m)
			if m.message["to"] == SERVER {
				user_connections[m.message["from"]] = m.sender_conn
			} else {
				if v, ok := user_connections[m.message["to"]]; ok {
					go func() {
						enc := gob.NewEncoder(v)
						enc.Encode(m.message)
					}()
				}
			}
		}
	}	
}
