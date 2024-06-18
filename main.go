package main

import (
	"log"
	"net/http"
)

type healthzHandler struct{}

func (hh healthzHandler) ServeHTTP(w http.ResponseWriter,_ *http.Request){
    w.WriteHeader(http.StatusOK)
    w.Header().Add("Content","text/plain")
    w.Header().Add("charset","utf-8")
    w.Write([]byte("OK"))
}

func main(){
    mux := http.NewServeMux()
    Server := &http.Server{
        Handler: mux,
        Addr: "localhost:8080",
    }
    fileSeverRoot := "."
    fileServerHander := http.FileServer(http.Dir(fileSeverRoot))
    mux.Handle("/app/*",http.StripPrefix("/app/", fileServerHander))
    mux.Handle("/healthz",healthzHandler{})
    log.Println("Listening on 8080")
    log.Fatal(Server.ListenAndServe())
}
