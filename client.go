package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
)

func main() {
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalf("Load keys %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", "127.0.0.1:8000", &config)
	if err != nil {
		log.Fatalf("Client: Dial %s", err)
	}
	defer conn.Close()
	log.Println("Client: Connected to ", conn.RemoteAddr())

	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		fmt.Println(v.Subject)
	}
	log.Println("Client: Handshake ", state.HandshakeComplete)
	log.Println("Client: Mutual ", state.NegotiatedProtocolIsMutual)

	message := "Hello\n"
	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalf("Client: write: %s", err)
	}
	log.Printf("Client: wrote %q (%d bytes)", message, n)

	reply := make([]byte, 256)
	n, err = conn.Read(reply)
	log.Printf("Client: read %q (%d bytes)", string(reply[:n]), n)
	log.Print("Client: exiting")
}
