package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"okv/client"

	"github.com/awesome-cap/hashmap"
)

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 10090, "The port to listen on; default is 8000.")
var isClient = flag.Bool("client", false, "The port to listen on; default is 8000.")

func main() {
	flag.Parse()
	if *isClient {
		client.Client()
		return
	}
	fmt.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s.\n", src)

	defer listener.Close()
	store := hashmap.New()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		go handleConnection(conn, store)
	}
}

func handleConnection(conn net.Conn, store *hashmap.HashMap) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)
	scanner := bufio.NewScanner(conn)
	for {
		ok := scanner.Scan()
		if !ok {
			break
		}
		handleMessage(scanner.Text(), conn, store)
	}

	fmt.Println("Client at " + remoteAddr + " disconnected.")
}

func handleMessage(message string, conn net.Conn, store *hashmap.HashMap) {
	// fmt.Println("< " + message)
	words := strings.Split(message, " ")
	// fmt.Printf("%v", words)
	if len(words) > 0 {
		switch words[0] {
		case "GET":
			fmt.Println("g")
			v, _ := store.Get(words[1])
			fmt.Printf("%v", v)
			conn.Write([]byte(v.(string) + "\n"))

		case "SET":
			fmt.Println("s")
			v := store.Set(words[1], words[2])
			if v == nil {
				conn.Write([]byte("OK\n"))
			} else {
				conn.Write([]byte(v.(string) + "\n"))
			}

		case "SHUTDOWN":
			os.Exit(0)

		default:
			conn.Write([]byte("Unrecognized command.\n"))
		}
	}
}
