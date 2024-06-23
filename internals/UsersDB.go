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

type user struct{
    Id int `json:"id"`
    Email string `json:"email"`
}

func NewUsersDB() (*db,error){
    if connected == true{
        log.Printf("There already exists a DB Connection. Be sure to close that before attempting to open a new one\n")
        return nil,errors.New("There already exists a DB Connection. Be sure to close that before attempting to open a new one")
    }
    connected = true
    database := db{mu : &sync.Mutex{},
                   db_path: "usersDatabase.json"}
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

func (database *db) AddUser(email string) (user,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        log.Printf("Error while loading database: %v\n",err)
        return user{},err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    users := []user{}
    data := struct{Mapper map[int]user `json:"users"`}{Mapper : make(map[int]user)}
    if len(json_data) != 0{
        if err := json.Unmarshal(json_data,&data); err != nil{
            log.Printf("Error while unmarshalling data : %v\n",err)
            return user{},err
        }
    }
    for _,val := range data.Mapper{
        users = append(users, val)
    }
    new_user := user{Id: len(users) + 1,
                     Email: email}
    users = append(users, user{Id: len(users) + 1,
                                  Email: email})
    data.Mapper[len(data.Mapper) + 1] = users[len(users) - 1]
    json_data,err = json.Marshal(data)
    if err != nil{
        log.Printf("Error while marshalling json: %v\n",err)
    }
    if err := os.WriteFile(database.db_path,json_data,0666); err != nil{
        log.Printf("Error while writing to database file: %v\n",err)
        return user{},nil
    }
    return new_user,err
}

func (database *db) QueryUsers() ([]user,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        log.Printf("Error while loading database: %v\n",err)
        return nil,err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    users := []user{}
    data := struct{Mapper map[int]user `json:"users"`}{Mapper : make(map[int]user)}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("Error while unmarshalling data : %v\n",err)
    }
    for _,val := range data.Mapper{
        users = append(users, val)
    }
    slices.SortFunc(users,func (a,b user) int{
                        if a.Id < b.Id{
                            return -1}else{
                                    return 1}})
    return users,nil
}

func (database *db) QueryUserByID(ID int) (user,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        log.Printf("Error while loading database: %v\n",err)
        return user{},err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    data := struct{Mapper map[int]user `json:"users"`}{Mapper : make(map[int]user)}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("Error while unmarshalling data : %v\n",err)
        return user{},err
    }
    for _,user := range data.Mapper{
        if user.Id == ID{
            return user,nil
        }
    }
    return user{},errors.New("Chirp Not Found")
}

func (database *db) CloseUsersDatabase() error{
    connected = false
    database = nil
    return nil
}
