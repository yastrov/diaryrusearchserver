package main

import (
	_ "encoding/json"
	"flag"
	"fmt"
	"github.com/yastrov/diaryrusearchserver/diaryruapi"
	"html/template"
	"log"
	"net/http"
	"time"
)

func HandlerAuthToDiary(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("assets/login.html")
		t.Execute(w, nil)
	} else { // POST
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form values!", http.StatusInternalServerError)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		sid, err := diaryruapi.Auth(username, password)
		if err != nil {
			http.Error(w, "Error Auth to diary.ru!", http.StatusInternalServerError)
			return
		}
		cookie := &http.Cookie{Name: "sid", Value: sid, Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: true}
		http.SetCookie(w, cookie)
		//fmt.Fprintf(w, sid)
		t, _ := template.ParseFiles("assets/search.html")
		t.Execute(w, sid)
	}
}

func HandlerSearch(w http.ResponseWriter, r *http.Request) {
	var err error
	var sid string
	if r.Method == "GET" {
		t, _ := template.ParseFiles("assets/search.html")
		cookie, err := r.Cookie("sid")
		if err != nil {
			sid = cookie.Value
			t.Execute(w, sid)
		} else {
			t.Execute(w, nil)
		}

	} else {
		if err = r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form values!", http.StatusInternalServerError)
			return
		}
		cookie, err := r.Cookie("sid")
		if err != nil {
			log.Printf(err.Error())
		}
		sid := cookie.Value
		fmt.Fprintf(w, sid)
	}
}

func main() {
	var hostport string
	portPtr := flag.String("port", ":8080", "port for listening")
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
		http.Redirect(w, r, "/login", 301)
	})
	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		log.Fatal("Stopped by user")
	})
	http.Handle("/static", http.StripPrefix("/static", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/login", HandlerAuthToDiary)
	http.HandleFunc("/search", HandlerSearch)
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
