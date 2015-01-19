package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func AuthToDiary(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("auth.html")
		t.Execute(w, p)

	} else { // POST
		username := r.FormValue("username")
		password := r.FormValue("password")
	}
}

func main() {
	var hostport string
	portPtr := flag.String("port", "10000", "port for listening")
	hostportPtr := flag.String("host", "", "host and port for listening")
	tlsPtr := flag.Bool("tls", false, "tls enabled (need cer.pem and key.pem files)")
	flag.Parse()
	if *hostportPtr != "" {
		hostport = *hostportPtr
	} else {
		hostport = *portPtr
	}
	// Handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", 0)
	})
	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		log.Fatal("Stopped by user")
	})
	http.HandleFunc("/login", AuthToDiary)
	// Serve
	if *tlsPtr == true {
		err := http.ListenAndServeTLS(hostport, "cert.pem", "key.pem", nil)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		if err := http.ListenAndServe(hostport, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
