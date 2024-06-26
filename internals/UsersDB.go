package internals

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"sync"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type user struct{
    Id int `json:"ID"`
    Email string `json:"Email"`
    Password string `json:"Password"`
    RefreshToken string `json:"refresh_token,omitempty"`
    RefreshTokenCreationTime time.Time `json:"refresh_token_creation_time,omitempty"`
}

type RequestUserInfo struct{Email string `json:"email"`
                            Password string `json:"password"`
                            ExpiresInSeconds time.Duration `json:"expires_in_seconds,omitempty"`}

type ResponseUserInfo struct{ID int `json:"id"`
                             Email string `json:"email"`
                             AccessToken string `json:"token,omitempty"`
                             RefreshToken string `json:"refresh_token,omitempty"`}

func NewUsersDB() (*db,error){
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

func (database *db) AddUser(userInfo RequestUserInfo) (ResponseUserInfo,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        log.Printf("Error while loading database: %v\n",err)
        return ResponseUserInfo{},err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    users := []user{}
    data := struct{Mapper map[string]user `json:"users"`}{Mapper : make(map[string]user)}
    if len(json_data) != 0{
        if err := json.Unmarshal(json_data,&data); err != nil{
            log.Printf("Error while unmarshalling data : %v\n",err)
            return ResponseUserInfo{},err
        }
    }
    for _,val := range data.Mapper{
        users = append(users, val)
    }
    if _,ok := data.Mapper[userInfo.Email];ok == true{
        return ResponseUserInfo{},errors.New("User with this email ID already exists")
    }
    respose := ResponseUserInfo{ID: len(users) + 1,
                                Email: userInfo.Email}
    hashed_password,err := bcrypt.GenerateFromPassword([]byte(userInfo.Password),10)
    users = append(users, user{Id: len(users) + 1,
                   Email: userInfo.Email,
                   Password : string(hashed_password)})
    data.Mapper[userInfo.Email] = users[len(users) - 1]
    json_data,err = json.Marshal(data)
    if err != nil{
        log.Printf("Error while marshalling json: %v\n",err)
    }
    if err := os.WriteFile(database.db_path,json_data,0666); err != nil{
        log.Printf("Error while writing to database file: %v\n",err)
        return ResponseUserInfo{},nil
    }
    return respose,err
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

func (database *db) QueryUserByID(ID int) (*user,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        log.Printf("Error while loading database: %v\n",err)
        return &user{},err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    data := struct{Mapper map[string]user `json:"users"`}{Mapper : make(map[string]user)}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("Error while unmarshalling data : %v\n",err)
        return &user{},err
    }
    for _,user := range data.Mapper{
        if user.Id == ID{
            return &user,nil
        }
    }
    return &user{},errors.New("Chirp Not Found")
}

func (database *db) UpdateUser(ID int, userInfo RequestUserInfo) error{
    u,err := database.QueryUserByID(ID)
    if err != nil{
        return err
    }
    json_data,err := database.loadDatabase()
    data := struct{Mapper map[string]user `json:"users"`}{Mapper : make(map[string]user)}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("Error while unmarshalling data : %v\n",err)
        return err
    }
    temp := data.Mapper[u.Email]
    delete(data.Mapper,u.Email)
    hashed_password,err := bcrypt.GenerateFromPassword([]byte(userInfo.Password),10)
    data.Mapper[userInfo.Email] = user{Email: userInfo.Email,
                                       Password: string(hashed_password),
                                       Id : u.Id,
                                       RefreshToken: temp.RefreshToken,
                                       RefreshTokenCreationTime: temp.RefreshTokenCreationTime}
    json_data,err = json.Marshal(&data)
    if err != nil{
        log.Printf("Error while marshalling data: %v",err)
        return err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    if err := os.WriteFile(database.db_path,json_data,0666); err != nil{
        log.Printf("Error while writing to database file: %v\n",err)
        return err
    }
    return nil
}

func (database *db) Login(userInfo RequestUserInfo) (ResponseUserInfo,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        log.Printf("Error while loading database : %v\n",err)
        return ResponseUserInfo{},err
    }
    data := struct{Mapper map[string]user `json:"users"`}{Mapper: make(map[string]user)}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("Error while unmarshalling JSON : %v\n",err)
        return ResponseUserInfo{},err
    }
    loginUser,ok := data.Mapper[userInfo.Email]
    if !ok{
        return ResponseUserInfo{},errors.New("User Email Not Found")
    }
    if bcrypt.CompareHashAndPassword([]byte(loginUser.Password),[]byte(userInfo.Password)) != nil{
        return ResponseUserInfo{},errors.New("Incorrect Password")
    }
    return ResponseUserInfo{ID: loginUser.Id,Email: loginUser.Email},nil
}

func (database *db) UpdateRefreshToken(Id int,refreshToken string) error{
    json_data,err := database.loadDatabase()
    if err != nil{
        fmt.Printf("Error while loading database : %v\n",err)
        return err
    }
    data := struct{Mapper map[string]*user `json:"users"`}{}
    if err := json.Unmarshal(json_data,&data); err != nil{
        fmt.Printf("Error while unmarshalling JSON : %v\n",err)
        return err
    }
    for _,user := range data.Mapper{
        if user.Id == Id{
            user.RefreshToken = refreshToken
            user.RefreshTokenCreationTime = time.Now()
        }
    }
    json_data,err = json.Marshal(data)
    if err != nil{
        fmt.Printf("Error while marshalling the data :%v",err)
        return err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    if err := os.WriteFile(database.db_path,json_data,0666); err != nil{
        log.Printf("Error while writing to database file: %v\n",err)
        return err
    }
    return nil
}

func (database *db) ValidateRefreshToken(refreshToken string) (*user,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        fmt.Printf("Error while loading database : %v\n",err)
        return &user{},err
    }
    data := struct{Mapper map[string]*user `json:"users"`}{}
    if err := json.Unmarshal(json_data,&data); err != nil{
        fmt.Printf("Error while unmarshalling JSON : %v\n",err)
        return &user{},err
    }
    for _,u := range data.Mapper{
        if u.RefreshToken == refreshToken{
            if time.Now().Sub(u.RefreshTokenCreationTime) > time.Hour * 24 * 60{
                return &user{},errors.New("Refresh Token Expired")
            }
            return u,nil
        }
    }
    return &user{},errors.New("Refresh Token Not Found")
}

func (database *db) DeleteRefreshToken(refreshToken string) error{
    json_data,err := database.loadDatabase()
    if err != nil{
        fmt.Printf("Error while loading database : %v\n",err)
        return err
    }
    data := struct{Mapper map[string]*user `json:"users"`}{}
    if err := json.Unmarshal(json_data,&data); err != nil{
        fmt.Printf("Error while unmarshalling JSON : %v\n",err)
        return err
    }
    deleted := false
    for _,u := range data.Mapper{
        if u.RefreshToken == refreshToken{
            u.RefreshToken = ""
            deleted = true
        }
    }
    if !deleted{
        return errors.New("Refresh Token Not Found")
    }
    json_data,err = json.Marshal(&data)
    if err != nil{
        fmt.Printf("Error while marshalling the data :%v",err)
        return err
    }
    database.mu.Lock()
    defer database.mu.Unlock()
    if err := os.WriteFile(database.db_path,json_data,0666); err != nil{
        log.Printf("Error while writing to database file: %v\n",err)
        return err
    }
    return nil
}
