// Simple HTTPS client using basic authentication.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"bufio"
	"crypto/tls"
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

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	f, err := os.Open(*filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	uri := "https://" + (*addr) + (*filename)
	fmt.Println(uri)

	// Set up HTTPS request with basic authorization.
	if *exec == "POST" || *exec == "PUT" {
		req, err := http.NewRequest("POST", uri, reader)
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.SetBasicAuth(*user, *pass)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HTTP Status:", resp.Status)
	fmt.Println("Response body:", string(html))
}
