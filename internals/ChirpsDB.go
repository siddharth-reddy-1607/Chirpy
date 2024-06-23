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

var connected bool = false

type db struct{
    mu *sync.Mutex
    db_path string
}

type chirp struct{
    Id int `json:"id"`
    Body string `json:"body"`
}

func NewDB() (*db,error){
    if connected == true{
        log.Printf("There already exists a DB Connection. Be sure to close that before attempting to open a new one\n")
        return nil,errors.New("There already exists a DB Connection. Be sure to close that before attempting to open a new one")
    }
    connected = true
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

func (database *db) Add(body string) (chirp,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        log.Printf("Error while loading database: %v\n",err)
        return chirp{},err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    chirps := []chirp{}
    data := struct{Mapper map[int]chirp `json:"chirps"`}{Mapper : make(map[int]chirp)}
    if len(json_data) != 0{
        if err := json.Unmarshal(json_data,&data); err != nil{
            log.Printf("Error while unmarshalling data : %v\n",err)
            return chirp{},err
        }
    }
    for _,val := range data.Mapper{
        chirps = append(chirps, val)
    }
    new_chirp := chirp{Id: len(chirps) + 1,
                       Body: body}
    chirps = append(chirps, chirp{Id: len(chirps) + 1,
                                  Body: body})
    data.Mapper[len(data.Mapper) + 1] = chirps[len(chirps) - 1]
    json_data,err = json.Marshal(data)
    if err != nil{
        log.Printf("Error while marshalling json: %v\n",err)
    }
    if err := os.WriteFile(database.db_path,json_data,0666); err != nil{
        log.Printf("Error while writing to database file: %v\n",err)
        return chirp{},nil
    }
    return new_chirp,err
}

func (database *db) Query() ([]chirp,error){
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
    return chirps,nil
}

func (database *db) QueryChirpByID(ID int) (chirp,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        log.Printf("Error while loading database: %v\n",err)
        return chirp{},err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    data := struct{Mapper map[int]chirp `json:"chirps"`}{Mapper : make(map[int]chirp)}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("Error while unmarshalling data : %v\n",err)
        return chirp{},err
    }
    for _,chirp := range data.Mapper{
        if chirp.Id == ID{
            return chirp,nil
        }
    }
    return chirp{},errors.New("Chirp Not Found")
}

func (database *db) CloseDatabase() error{
    connected = false
    database = nil
    return nil
}
