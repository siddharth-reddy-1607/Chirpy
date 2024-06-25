package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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
        database,err := internals.NewChirpsDB()
        if err != nil{
            err_msg := fmt.Sprintf("Error while creating new DB Connection : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        responseJSON,err := database.AddChirp(requestJSON.Body)
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
        database,err := internals.NewChirpsDB()
        if err != nil{
            err_msg := fmt.Sprintf("Error while creating new DB Connection : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        chirps,err := database.QueryChirps()
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

func postUsersHandler() http.Handler{
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
        decoder := json.NewDecoder(r.Body)
        requestJSON := struct{Email string `json:"email"`;
                              Password string `json:"password"`}{}
        if err := decoder.Decode(&requestJSON); err != nil{
            err_msg := fmt.Sprintf("Error while decoding json :%v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        database,err := internals.NewUsersDB()
        if err != nil{
            err_msg := fmt.Sprintf("Error while creating new DB Connection : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        responseJSON,err := database.AddUser(internals.RequestUserInfo{Email: requestJSON.Email,
                                                                       Password: requestJSON.Password})
        if err!= nil{
            err_msg := fmt.Sprintf("Error while adding chirp : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        enconder := json.NewEncoder(w)
        w.WriteHeader(http.StatusCreated)
        enconder.Encode(responseJSON)
    })
}

func getChirpByIDHandler() http.Handler{
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
        id,err := strconv.Atoi(r.PathValue("chirpID"))
        if err!= nil{
            err_msg := fmt.Sprintf("chirpID must be an integer : %v",err)
            http.Error(w,err_msg,http.StatusNotFound)
            return
        }
        db,err := internals.NewChirpsDB()
        if err != nil{
            err_msg := fmt.Sprintf("Error while connecting to the DB : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        chirp,err := db.QueryChirpByID(id)
        if err != nil{
            if err.Error() == "Chirp Not Found"{
                http.Error(w,err.Error(),http.StatusNotFound)
                return
            } 
            err_msg := fmt.Sprintf("Error while getting chirp with ID %d : %v",id,err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        encoder := json.NewEncoder(w)
        if err := encoder.Encode(chirp); err != nil{
            err_msg := fmt.Sprintf("Error while encoding JSON : %v",err)
            http.Error(w,err_msg,http.StatusNotFound)
            return
        }
    })
}

func LoginHandler() http.Handler{
    return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request){
        requestJSON := internals.RequestUserInfo{ExpiresInSeconds: -1}
        decoder := json.NewDecoder(r.Body)
        if err := decoder.Decode(&requestJSON); err != nil{
            err_msg := fmt.Sprintf("Error while decoding json : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        database,err := internals.NewUsersDB()
        if err != nil{
            err_msg := fmt.Sprintf("Error while creating database conenction : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        responseJSON,err := database.Login(requestJSON)
        if err != nil{
            if err.Error() == "Incorrect Password"{
                w.WriteHeader(http.StatusUnauthorized)
                w.Write([]byte("Incorrect Password"))
                return
            }
            err_msg := fmt.Sprintf("Error occured while logging is user: %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return

        }
        if requestJSON.ExpiresInSeconds == -1{
            requestJSON.ExpiresInSeconds = time.Second*60*60*24
        }else{
            requestJSON.ExpiresInSeconds *= time.Second
        }
        token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.RegisteredClaims{Issuer: "chirpy",
                                                                               IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
                                                                               ExpiresAt: jwt.NewNumericDate(time.Now().Add(requestJSON.ExpiresInSeconds)),
                                                                               Subject: strconv.Itoa(responseJSON.ID)})
       signed_token,err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
        if err != nil{
            err_msg := fmt.Sprintf("Error will signing in the token : %v",err)
            http.Error(w,err_msg,http.StatusUnauthorized)
            return
        } 
        responseJSON.Signed_Token = signed_token
        encoder := json.NewEncoder(w)
        encoder.Encode(&responseJSON)
    })
}

func UpdateUserDetailsHandler() http.Handler{
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
        auth := r.Header.Get("Authorization")
        if auth == ""{
            err_msg := "Please login before trying to update your details"
            http.Error(w,err_msg,http.StatusBadRequest)
            return
        }
        token_string,ok := strings.CutPrefix(auth,"Bearer ")
        fmt.Printf("Token is %v\n",token_string)
        if !ok {
            err_msg := "Invalid Header. Missing Prefix 'Bearer '"
            http.Error(w,err_msg,http.StatusBadRequest)
            return
        }
        claims := jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(token_string, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_KEY")), nil
        })
        if err != nil{
            err_msg := fmt.Sprintf("Error parsing the token string : %v :",err)
            http.Error(w,err_msg,http.StatusUnauthorized)
            return
        }
        subject,err := token.Claims.GetSubject()
        if err != nil{
            err_msg := fmt.Sprintf("Error while getting subject : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        ID,err := strconv.Atoi(subject)
        fmt.Printf("ID of user to update is %d\n",ID)
        if err != nil{
            err_msg := fmt.Sprintf("Error while converting token to integer : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        database,err := internals.NewUsersDB()
        if err != nil{
            err_msg := fmt.Sprintf("Error while creating database conenction : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        request_json := internals.RequestUserInfo{}
        decoder := json.NewDecoder(r.Body)
        decoder.Decode(&request_json)
        if err := database.UpdateUser(ID,request_json); err != nil{
            err_msg := fmt.Sprintf("Error while updating users email and password : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }

        user,err := database.QueryUserByID(ID)

        responseJSON := internals.ResponseUserInfo{Email: user.Email,
                                                   ID: user.Id}
        if err != nil{
            err_msg := fmt.Sprintf("Error while finding user with ID : %v",err)
            http.Error(w,err_msg,http.StatusInternalServerError)
            return
        }
        encoder := json.NewEncoder(w)
        encoder.Encode(&responseJSON)
     })
}

func main(){
    debug := flag.Bool("debug",false,"Debug mode")
    flag.Parse()
    if err := godotenv.Load(); err != nil{
        log.Fatalf("Error while loading .env file : %v\n",err)
    }
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
    mux.Handle("GET /api/chirps/{chirpID}",getChirpByIDHandler())
    mux.Handle("POST /api/chirps",postChirpsHandler())
    mux.Handle("POST /api/users",postUsersHandler())
    mux.Handle("POST /api/login",LoginHandler())
    mux.Handle("PUT /api/users",UpdateUserDetailsHandler())
    log.Println("Listening on 8080")
    log.Fatal(Server.ListenAndServe())
    if *debug == true{
        fmt.Printf("Removing Database files\n")
        os.Remove("usersDatabase.json")
        os.Remove("chirpsDatabase.json")
    }
}
