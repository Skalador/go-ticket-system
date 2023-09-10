package main

import (
	"io"
	"log"
	"net/http"
	"text/template"
)

const PORT = ":8000"

func helloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	t,err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w,nil)
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //Fileserver to load css
	http.HandleFunc("/",rootHandler)
	http.HandleFunc("/hello", helloHandler) //Testing purpose
    log.Println("Listing for requests at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(PORT, nil))
}
