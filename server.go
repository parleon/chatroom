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

	// goroutine to accept incoming connections and route messages into message queue
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
					decerr := dec.Decode(&message_map)
					if decerr != nil {
						fmt.Println("decerr")
						conn.Close()
						return
					} else {
					m_queue <- message_wrapper{message: message_map, sender_conn: conn}
					}
				}
			}()
		}
	}()

	// process message queue
	user_connections := make(map[string]net.Conn)
	connection_encodings := make(map[string]*gob.Encoder)
	for {
		select {
		case m := <-m_queue:
			fmt.Println(m)
			if m.message["to"] == SERVER {
				fmt.Println("Registering User")
				user_connections[m.message["from"]] = m.sender_conn
				connection_encodings[m.message["from"]]= gob.NewEncoder(m.sender_conn)
			} else {
				if v, ok := user_connections[m.message["to"]]; ok {
					go func() {
						encerr := connection_encodings[m.message["to"]].Encode(m.message)
						if encerr != nil {
							fmt.Println("encerr")
							delete(user_connections,m.message["to"])
							v.Close()
						}
					}()
				}
			}
		}
	}	
}
