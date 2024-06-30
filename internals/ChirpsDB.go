package internals

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"sync"
)

type db struct{
    mu *sync.Mutex
    db_path string
}

type chirp struct{
    Id int `json:"id"`
    AuthorId int `json:"author_id"`
    Body string `json:"body"`
}

type RequestChirpInfo struct{
    Author_Id int `json:"author_id"`
    Body string `json:"body"`
}

type ResponseChirpInfo struct{
    Id int `json:"id"`
    AuthorId int `json:"author_id"`
    Body string `json:"body"`
}

func NewChirpsDB() (*db,error){
    database := db{mu : &sync.Mutex{},
                   db_path: "chirpsDatabase.json"}
   _,err := os.Stat(database.db_path)
   if os.IsNotExist(err) == true{
       fmt.Printf("Database file not present. Creating it\n")
       if _,err := os.Create(database.db_path); err != nil{
           log.Printf("Error while creating database : %v\n",err)
           return nil,err
       }
   }
    return &database,nil
}

func (database *db) loadDatabase() ([]byte,error){
    database.mu.Lock()
    defer database.mu.Unlock()
    json_data,err := os.ReadFile(database.db_path)
    if err != nil {
        log.Printf("Error while reading database file: %v\n",err)
        return nil,err
    }
    return json_data,nil
}

func (database *db) AddChirp(chirpInfo RequestChirpInfo) (ResponseChirpInfo,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        log.Printf("Error while loading database: %v\n",err)
        return ResponseChirpInfo{},err
    }
    chirps := []chirp{}
    data := struct{Mapper map[int]chirp `json:"chirps"`}{Mapper : make(map[int]chirp)}
    if len(json_data) != 0{
        if err := json.Unmarshal(json_data,&data); err != nil{
            log.Printf("Error while unmarshalling data : %v\n",err)
            return ResponseChirpInfo{},err
        }
    }
    for _,val := range data.Mapper{
        chirps = append(chirps, val)
    }
    new_chirp := chirp{Id: len(chirps) + 1,
                       AuthorId: chirpInfo.Author_Id,
                       Body: chirpInfo.Body}
    chirps = append(chirps, new_chirp)
    data.Mapper[len(data.Mapper) + 1] = chirps[len(chirps) - 1]
    json_data,err = json.Marshal(&data)
    if err != nil{
        log.Printf("Error while marshalling json: %v\n",err)
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    if err := os.WriteFile(database.db_path,json_data,0666); err != nil{
        log.Printf("Error while writing to database file: %v\n",err)
        return ResponseChirpInfo{},nil
    }
    return ResponseChirpInfo{Id:new_chirp.Id,Body:new_chirp.Body,AuthorId: chirpInfo.Author_Id},err
}

func (database *db) QueryChirps() ([]ResponseChirpInfo,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        log.Printf("Error while loading database: %v\n",err)
        return nil,err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    chirps := []chirp{}
    data := struct{Mapper map[int]chirp `json:"chirps"`}{Mapper : make(map[int]chirp)}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("Error while unmarshalling data : %v\n",err)
    }
    for _,val := range data.Mapper{
        chirps = append(chirps, val)
    }
    slices.SortFunc(chirps,func (a,b chirp) int{
                        if a.Id < b.Id{
                            return -1}else{
                                    return 1}})
    responseChirps := []ResponseChirpInfo{}
    for _,c := range chirps{
        responseChirp := ResponseChirpInfo{Id : c.Id,
                                           Body: c.Body,
                                           AuthorId: c.AuthorId}
        responseChirps = append(responseChirps,responseChirp)
    }
    return responseChirps,nil
}

func (database *db) QueryChirpByID(ID int) (ResponseChirpInfo,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        log.Printf("Error while loading database: %v\n",err)
        return ResponseChirpInfo{},err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    data := struct{Mapper map[int]chirp `json:"chirps"`}{Mapper : make(map[int]chirp)}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("Error while unmarshalling data : %v\n",err)
        return ResponseChirpInfo{},err
    }
    for _,chirp := range data.Mapper{
        if chirp.Id == ID{
            responseChirp :=  ResponseChirpInfo{Id : chirp.Id,
                              Body: chirp.Body,
                              AuthorId: chirp.AuthorId}
            return responseChirp,nil
        }
    }
    return ResponseChirpInfo{},errors.New("Chirp Not Found")
}

func (database *db) DeleteChirp(chirpID int, authorID int) error{
    json_data,err := database.loadDatabase()
    if err != nil{
        log.Printf("Error while loading database: %v\n",err)
        return err
    }
    data := struct{Mapper map[int]chirp `json:"chirps"`}{Mapper : make(map[int]chirp)}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("Error while unmarshalling data : %v\n",err)
        return err
    }
    found := false
    for _,c := range data.Mapper{
        if c.Id == chirpID{
            if c.AuthorId != authorID{
                return errors.New("Forbidden")
            }
            found = true
            break
        }
    }
    if !found{
        return fmt.Errorf("Chirp with ID %d not found",chirpID)
    }
    delete(data.Mapper,chirpID)
    return nil
}
