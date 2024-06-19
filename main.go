package main

import (
	"fmt"
	"log"
	"net/http"
)


type healthzHandler struct{}

func (hh *healthzHandler) ServeHTTP(w http.ResponseWriter,_ *http.Request){
    w.WriteHeader(http.StatusOK)
    w.Header().Add("Content","text/plain")
    w.Header().Add("charset","utf-8")
    w.Write([]byte("OK"))
}

type PageMetrics struct{
    timesVisted int
}
func (pm *PageMetrics) middlewareUpdateHomePageMetrics (next http.Handler) http.Handler{
    icr := func(w http.ResponseWriter, r *http.Request){
        pm.timesVisted += 1
        next.ServeHTTP(w,r)
    }
    return http.HandlerFunc(icr)
}
func (pm *PageMetrics) middlewareResetHomePageMetrics(next http.Handler) http.Handler{
    reset := func(w http.ResponseWriter,r *http.Request){
        pm.timesVisted = 0
        next.ServeHTTP(w,r)
    }
    return http.HandlerFunc(reset)
}
func (pm *PageMetrics) displayMetrics() http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        hits := fmt.Sprintf("Hits: %d",pm.timesVisted)
        w.Write([]byte(hits))
    })
}
func (pm *PageMetrics) resetMetrics() http.Handler{
    return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request){
        pm.timesVisted = 0
    })
}

func addFSHeader(h http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        w.Header().Add("Cache-Control","no-cache")
        h.ServeHTTP(w,r)
    })
}

func main(){
    mux := http.NewServeMux()
    Server := &http.Server{
        Handler: mux,
        Addr: "localhost:8080",
    }
    fileSeverRoot := "."
    fileServerHander := http.FileServer(http.Dir(fileSeverRoot))
    pm := PageMetrics{timesVisted : 0}
    mux.Handle("/app/*",addFSHeader(pm.middlewareUpdateHomePageMetrics(http.StripPrefix("/app", fileServerHander))))
    mux.Handle("GET /healthz",&healthzHandler{})
    mux.Handle("GET /metrics",pm.displayMetrics())
    mux.Handle("/reset",pm.resetMetrics())
    log.Println("Listening on 8080")
    log.Fatal(Server.ListenAndServe())
}
