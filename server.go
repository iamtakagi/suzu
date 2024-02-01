package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
)

func main() {
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatalf("Load keys %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	service := "0.0.0.0:8000"
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("Server: listen: %s", err)
	}
	log.Print("Server: listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Server: accept: %s", err)
			break
		}
		defer conn.Close()
		log.Printf("Server: Accepted from %s", conn.RemoteAddr())
		tlscon, ok := conn.(*tls.Conn)
		if ok {
			log.Print("ok=true")
			state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
		}
		go HandleRequest(conn)
	}
}

func HandleRequest(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		log.Print("Server: Waiting")
		n, err := conn.Read(buf)
		if err != nil {
			if err != nil {
				log.Printf("Server: Read %s", err)
			}
			break
		}
		log.Printf("Server: Echo %q\n", string(buf[:n]))
		n, err = conn.Write(buf[:n])

		n, err = conn.Write(buf[:n])
		log.Printf("Server: Wrote %d bytes", n)

		if err != nil {
			log.Printf("Server: Write %s", err)
			break
		}
	}
	log.Println("Server: closed")
}
