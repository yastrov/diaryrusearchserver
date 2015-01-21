package main

import (
	_ "encoding/json"
	"flag"
	_ "fmt"
	"github.com/yastrov/diaryrusearchserver/diaryruapi"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

type MyResponse struct {
	SID      string
	Comments []*diaryruapi.CommentStruct
	Posts    []*diaryruapi.PostStruct
	//Umails []*diaryruapi.
}

func HandlerAuthToDiary(w http.ResponseWriter, r *http.Request) {
	my_response := &MyResponse{SID: "", Comments: make([]*diaryruapi.CommentStruct, 0), Posts: make([]*diaryruapi.PostStruct, 0)}
	if r.Method == "GET" {

		t, _ := template.ParseFiles("assets/login.html")
		t.Execute(w, my_response)
	} else { // POST
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			http.Error(w, "Error parsing form values!", http.StatusInternalServerError)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		sid, err := diaryruapi.Auth(username, password)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error Auth to diary.ru!", http.StatusInternalServerError)
			return
		}
		cookie := &http.Cookie{Name: "sid", Value: sid, Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: true}
		http.SetCookie(w, cookie)
		//fmt.Fprintf(w, sid)
		my_response.SID = sid
		t, err := template.ParseFiles("assets/search.html")
		if err != nil {
			log.Println(err)
		}
		t.Execute(w, my_response)
	}
}

func HandlerSearch(w http.ResponseWriter, r *http.Request) {
	var err error
	var sid string
	my_response := &MyResponse{SID: "", Comments: nil, Posts: nil}
	if r.Method == "GET" {
		t, _ := template.ParseFiles("assets/search.html")
		cookie, err := r.Cookie("sid")
		if err != nil {
			log.Println(err)
			sid = cookie.Value
			my_response.SID = sid
			t.Execute(w, my_response)
		} else {
			t.Execute(w, my_response)
		}

	} else {
		var wg sync.WaitGroup
		if err = r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form values!", http.StatusInternalServerError)
			return
		}
		var sid string
		cookie, err := r.Cookie("sid")
		if err != nil || cookie.Value == "" {
			log.Println(err)
			log.Printf(err.Error())
			sid = r.FormValue("sid")
		} else {
			sid = cookie.Value
		}
		diarytype := r.FormValue("diarytype")
		log.Println("diarytype: ", diarytype)
		//keyword := r.FormValue("keyword")
		shortname := r.FormValue("shortname")
		log.Println("shortname: ", shortname)
		switch diarytype {
		case "diary":
			log.Println("diary")
			post_chan := make(chan *diaryruapi.PostStruct)
			err_chan := make(chan error)
			journal, err := diaryruapi.JournalGet(sid, "", shortname)
			if err != nil {
				log.Panic(err)
			}
			wg.Add(1)
			go diaryruapi.PostsAllGetChannels(sid, diarytype, journal, post_chan, err_chan, &wg)
			wg.Wait()
			result := make([]*diaryruapi.PostStruct, 0, 20)
			var post *diaryruapi.PostStruct
			for {
				select {
				case post = <-post_chan:
					result = append(result, post)
				case err = <-err_chan:
					log.Println(err)
				}
			}
			my_response.Posts = result
		case "umail":
			log.Println("umail")
		default:
			log.Println("default")
		}
		t, _ := template.ParseFiles("assets/search.html")
		t.Execute(w, my_response)
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
