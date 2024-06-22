package main

import (
	"encoding/json"
    "errors"
	"fmt"
	"strings"
	"log"
	"net/http"
    "github.com/siddharth-reddy-1607/Chirpy/internals"
)

type ResponseData struct{
    WriteHeader int
    Header map[string]string
    Data any
}
type healthzHandler struct{}

func (hh *healthzHandler) ServeHTTP(w http.ResponseWriter,_ *http.Request){
    w.WriteHeader(http.StatusOK)
    w.Header().Add("Content","text/plain")
    w.Header().Add("charset","utf-8")
    w.Write([]byte("OK"))
}

func addFSHeader(h http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        w.Header().Add("Cache-Control","no-cache")
        h.ServeHTTP(w,r)
    })
}

func filter_profane_words(chirp *string) string{
    profaneWordMappings := map[string]bool{
        "kerfuffle" : true,
        "sharbert" : true,
        "fornax" : true,
    }
    words := strings.Split(*chirp, " ")
    filter_words := []string{}
    for _,word := range words{
        if _,ok := profaneWordMappings[strings.ToLower(word)]; ok{
            filter_words = append(filter_words,"****")
        }else{
            filter_words = append(filter_words,word)
        }
    }
    return strings.Join(filter_words," ")
}

func validate_chirp_length(body *string) (error,any){
    if len(*body) > 140{
        response := ResponseData{
            WriteHeader: http.StatusBadRequest,
            Data: struct{Valid bool `json:"valid"`}{Valid: false},
        }  
        return errors.New("Length of chirp body is greater than 140"),response
    }
    return nil,nil
}

func postChirpsHandler() http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        requestJSON := struct{Body string `json:"body"`}{}
        decoder := json.NewDecoder(r.Body)
        if err := decoder.Decode(&requestJSON); err != nil{
            err_msg := fmt.Sprintf("Error while decoding json here: %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        database,err := internals.NewDB()
        if err != nil{
            err_msg := fmt.Sprintf("Error while creating new DB Connection : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        defer database.CloseDatabase()
        responseJSON,err := database.Add(requestJSON.Body)
        if err!= nil{
            err_msg := fmt.Sprintf("Error while adding chirp : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
        }
        enconder := json.NewEncoder(w)
        w.WriteHeader(http.StatusCreated)
        enconder.Encode(responseJSON)
    })
}

func getChirpsHandler() http.Handler{
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
        database,err := internals.NewDB()
        if err != nil{
            err_msg := fmt.Sprintf("Error while creating new DB Connection : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        defer database.CloseDatabase()
        chirps,err := database.Query()
        if err!= nil{
            err_msg := fmt.Sprintf("Error while querying the database : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
        }
        encoder := json.NewEncoder(w)
        for _,chirp := range chirps{
            chirp.Body = filter_profane_words(&chirp.Body)
        }
        encoder.Encode(chirps)
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
    am := AppMetrics{TimesVisited : 0}
    mux.Handle("/app/*",addFSHeader(am.middlwareUpdateMetrics(http.StripPrefix("/app", fileServerHander))))
    mux.Handle("GET /api/healthz",&healthzHandler{})
    mux.Handle("GET /admin/metrics",am.displayMetrics())
    mux.Handle("/api/reset",am.resetMetrics())
    mux.Handle("GET /api/chirps",getChirpsHandler())
    mux.Handle("POST /api/chirps",postChirpsHandler())
    log.Println("Listening on 8080")
    log.Fatal(Server.ListenAndServe())
}
