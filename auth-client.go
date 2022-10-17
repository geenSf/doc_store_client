// Simple HTTPS client using basic authentication.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a arguments!")
		return
	}

	addr := flag.String("addr", "", "HTTPS server address")
	certFile := flag.String("certfile", "cert.pem", "trusted CA certificate")
	user := flag.String("user", "", "username")
	pass := flag.String("pass", "", "password")
	filename := flag.String("file", "", "JSON message file")
	exec := flag.String("exec", "", "execute REST method")
	key := flag.String("key", "", "set key for GET method")
	flag.Parse()

	if *addr == "" {
		log.Fatalf("Specify the request address!")
	}

	if *exec == "" {
		log.Fatalf("Specify the request type!")
	}

	if *exec == "POST" && *filename == "" {
		log.Fatalf("Specify the JSON message file!")
	}

	// Read the trusted CA certificate from a file and set up a client with TLS
	// config to trust a server signed with this certificate.
	cert, err := os.ReadFile(*certFile)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("unable to parse cert from %s", *certFile)
	}

	//create client
	client := NewClient(*addr, certPool)

	var res *http.Response

	// Set up HTTPS request with basic authorization.

	// POST request
	if *exec == "POST" || *exec == "PUT" {
		res, err := client.PostReq(*filename, *user, *pass)
		if err != nil {
			log.Fatal(err)
		}
	}

	// GET request
	if *exec == "GET" {

		res, err := client.GetReq(key)
		if err != nil {
			log.Fatal(err)
		}

	}

	html, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HTTP Status:", res.Status)
	fmt.Println("Response body:", string(html))
}
