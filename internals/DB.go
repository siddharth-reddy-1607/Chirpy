package internals

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
    "time"
)

var (
    dbNotPresentError = errors.New("Database file is not present")
    dbLoadError = errors.New("Error while loading the database")
    dbWriteError = errors.New("Error while writing to the database")
    DataUnmarshalError = errors.New("Error while unmarshalling data")
    DataMarshalError = errors.New("Error while marshalling data")
)

type db struct{
    mu *sync.Mutex
    dbPath string
}

type User struct{
    ID int `json:"id"`
    Email string `json:"email"`
    EncryptedPassword string `json:"encryptedPassword"`
    IsChirpyRed bool `json:"is_chirpy_red"`
}

type Chirp struct{
    ID int `json:"id"`
    Body string `json:"body"`
    AuthorID int `json:"author_id"`
}

type dbStructure struct{
    Users map[int]User `json:"users"`
    Chirps map[int]Chirp `json:"chirps"`
    RefreshToken map[string]struct{UserID int
                                   ExpirationTime time.Time} `json:"refreshToken"`
}

func NewDB(dbFileName string) (*db,error){
    database := &db{mu : &sync.Mutex{},
                  dbPath: dbFileName}
    err := database.doesDBExist()
    if err == nil{
        return database,nil
    }
    if err == dbNotPresentError{
        fmt.Printf("Database file not present. Creating it\n")
        if _,err := os.Create(database.dbPath); err != nil{
            log.Printf("Error while creating database: %v\n",err)
            return nil,err
        }
        return database,nil
    }
    log.Printf("Error while getting info database for DB File: %v\n",err)
    return nil,err
}

func(database *db) doesDBExist() error{
    if _,err := os.Stat(database.dbPath); err != nil{
        if os.IsNotExist(err) == true{
            return dbNotPresentError
        }
        return err
    }
    return nil
}

func (database *db) loadDatabase() ([]byte,error){
    if err := database.doesDBExist(); err != nil{
        log.Printf("Database file doesn't exist: %v\n",err)
        return nil,err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    json_data,err := os.ReadFile(database.dbPath)
    if err != nil {
        log.Printf("%v: %v\n",dbWriteError,err)
        return nil,dbWriteError
    }
    return json_data,nil
}

func (database *db) writeDatabase(json_data []byte) error{
    if err := database.doesDBExist(); err != nil{
        log.Printf("Database file doesn't exist: %v\n",err)
        return err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    if err := os.WriteFile(database.dbPath,json_data,0666); err != nil{
        log.Printf("%v: %v\n",dbWriteError,err) 
        return dbWriteError
    }
    return nil
}

func (database *db) EraseDB() error{
    err := database.doesDBExist()
    if err == nil{
        if err := os.Remove(database.dbPath); err != nil{
            log.Printf("Error while removing DB File: %v\n",err)
            return errors.New("Error while removing the database file")
        }
    }
    if err == dbNotPresentError{
        return nil
    }
    log.Printf("Error while getting info database for DB File: %v\n",err)
    return err
}
