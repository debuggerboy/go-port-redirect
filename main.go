package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

var (
	listen int
	target int
	tghost string
)

func handleConnection(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Println("Error connecting to the target:", err)
		return
	}

	defer targetConn.Close()

	go func() {
		_, err := io.Copy(targetConn, clientConn)
		if err != nil {
			log.Println("Error copying data from client to target:", err)
			return
		}
	}()
	_, err = io.Copy(clientConn, targetConn)
	if err != nil {
		log.Println("Error copying data from target to client:", err)
	}
}

func main() {
	flag.IntVar(&listen, "listen", 0, "Specify a listener port")
	flag.IntVar(&target, "target", 0, "Specify a target port")
	flag.StringVar(&tghost, "targethost", "localhost", "Specify a target host (optional)")
	flag.IntVar(&listen, "l", 0, "Specify a listener port")
	flag.IntVar(&target, "t", 0, "Specify a target port")
	flag.StringVar(&tghost, "h", "localhost", "Specify a target host (optional)")
	flag.Parse()

	filename := filepath.Base(os.Args[0])

	if listen == 0 {
		fmt.Println("Error: Listen port cannot be empty")
		fmt.Println("try", filename, "--help")
		os.Exit(3)
	}

	if target == 0 {
		fmt.Println("Error: Target port cannot be empty")
		fmt.Println("try ", filename, " --help")
		os.Exit(3)
	}

	listenPort := strconv.Itoa(listen)
	targetPort := strconv.Itoa(target)

	listenPortLiteral := ":" + listenPort
	targetPortLiteral := tghost + ":" + targetPort

	/*
		  * debugging port variables
			fmt.Printf("%s\n", listenPortLiteral)
			fmt.Printf("%s\n", targetPortLiteral)
			os.Exit(3)
	*/

	listener, err := net.Listen("tcp", listenPortLiteral)
	if err != nil {
		log.Fatal("Error while starting listener:", err)
	}
	defer listener.Close()

	log.Println("Listening on port", listenPort, ", and redirecting traffic to port", targetPortLiteral)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(clientConn, targetPortLiteral)
	}
}
