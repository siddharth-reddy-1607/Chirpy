package internals

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
)

var (
    UserNotFoundError = errors.New("User Not Found")
)
var newUserID int = 1

func (database *db) AddUser(email, encryptedPassword string) (User,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        return User{},err
    }
    data := dbStructure{}
    if len(json_data) != 0{
        if err := json.Unmarshal(json_data,&data); err != nil{
            log.Printf("%v: %v\n",DataUnmarshalError,err)
            return User{},DataUnmarshalError
        }
    }
    if data.Users == nil{
        data.Users = make(map[int]User)
    }
    data.Users[newUserID] = User{ID: newUserID,
                                 Email: email,
                                 IsChirpyRed: false,
                                 EncryptedPassword: encryptedPassword}
    json_data,err = json.Marshal(&data)
    if err != nil{
        log.Printf("%v: %v\n",DataMarshalError,err)
        return User{},DataMarshalError
    }
    if err := database.writeDatabase(json_data); err != nil{
        return User{},err
    }
    newUserID += 1
    return data.Users[newUserID - 1],nil
}

func (database *db) GetUsers() ([]User,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        return nil,err
    }
    data := dbStructure{}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("%v: %v\n",DataUnmarshalError,err)
        return nil,DataUnmarshalError
    }
    if data.Users == nil{
        data.Users = make(map[int]User)
    }
    users := []User{}
    for _,val := range data.Users{
        users = append(users, val)
    }
    slices.SortFunc(users,func (a,b User) int{
                            if a.ID < b.ID{
                                return -1
                            }else{
                                return 1
                            }
                          })
    return users,nil
}

func (database *db) GetUserByID(ID int) (User,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        return User{},err
    }
    data := dbStructure{}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("%v: %v\n",DataUnmarshalError,err)
        return User{},DataUnmarshalError
    }
    if data.Users == nil{
        data.Users = make(map[int]User)
    }
    if user,ok := data.Users[ID]; ok{
        return user,nil
    }
    return User{},UserNotFoundError
}

func (database *db) GetUserByEmail(email string) (User,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        return User{},err
    }
    data := dbStructure{}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("%v: %v\n",DataUnmarshalError,err)
        return User{},DataUnmarshalError
    }
    if data.Users == nil{
        data.Users = make(map[int]User)
    }
    for _,user := range data.Users{
        if user.Email == email{
            return user,nil
        }
    }
    return User{},UserNotFoundError
}

func (database *db) UpdateUser(ID int, email,encryptedPassword string) (User,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        return User{},err
    }
    data := dbStructure{}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("%v : %v\n",DataUnmarshalError,err)
        return User{},DataUnmarshalError
    }
    if data.Users == nil{
        data.Users = make(map[int]User)
    }
    user,ok := data.Users[ID]
    if !ok{
        return User{},UserNotFoundError
    }
    user.Email = email
    user.EncryptedPassword = encryptedPassword
    data.Users[ID] = User{ID : user.ID,
                          Email : user.Email, 
                          EncryptedPassword : user.EncryptedPassword,
                          IsChirpyRed : user.IsChirpyRed}
    fmt.Printf("Data after update %+v\n",data.Users)
    json_data,err = json.Marshal(&data)
    if err != nil{
        log.Printf("%v: %v\n",DataMarshalError,err)
        return User{},err
    }
    if err := database.writeDatabase(json_data); err != nil{
        return User{},err
    }
    return data.Users[ID],nil
}

func (database*db) UpgradeUser(ID int) (User,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        return User{},err
    }
    data := dbStructure{}
    if err := json.Unmarshal(json_data,&data); err != nil{
        fmt.Printf("%v: %v\n",DataUnmarshalError,err)
        return User{},err
    }
    if data.Users == nil{
        data.Users = make(map[int]User)
    }
    user,ok := data.Users[ID]
    if !ok{
        return User{},UserNotFoundError
    }
    user.IsChirpyRed = true
    data.Users[ID] = User{ID : user.ID,
                          Email : user.Email, 
                          EncryptedPassword : user.EncryptedPassword,
                          IsChirpyRed : user.IsChirpyRed}
    json_data,err = json.Marshal(&data)
    if err != nil{
        fmt.Printf("%v: %v",DataMarshalError,err)
        return User{},err
    }
    if err := database.writeDatabase(json_data); err != nil{
        return User{},err
    }
    return data.Users[ID],nil
}
