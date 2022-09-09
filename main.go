package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/imlonghao/ja3-server/crypto/tls"
	"github.com/imlonghao/ja3-server/net/http"
	"net"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	hash := md5.Sum([]byte(r.JA3Fingerprint))
	out := make([]byte, 32)
	hex.Encode(out, hash[:])

	fmt.Println(r.JA3Fingerprint)

	w.Write(out)
}

func main() {
	handler := http.HandlerFunc(handler)
	server := &http.Server{Addr: ":8443", Handler: handler}

	ln, err := net.Listen("tcp", ":8443")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}
	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}}

	tlsListener := tls.NewListener(ln, &tlsConfig)
	fmt.Println("HTTP up.")
	err = server.Serve(tlsListener)
	if err != nil {
		panic(err)
	}

	ln.Close()
}
