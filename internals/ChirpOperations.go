package internals

import (
	"encoding/json"
	"errors"
	"log"
    "slices"
)

var(
    ChirpNotFoundError = errors.New("Chirp Not Found")
)

var newChirpID int = 1

func(database *db) AddChirp(body string,authorID int) (Chirp,error){
    json_data,err := database.loadDatabase() 
    if err != nil{
        return Chirp{},err
    }
    data := dbStructure{}
    if len(json_data) != 0{
        if err := json.Unmarshal(json_data,&data); err != nil{
            log.Printf("%v: %v",DataUnmarshalError,err)
            return Chirp{},DataUnmarshalError
        }
    }
    if data.Chirps == nil{
        data.Chirps = make(map[int]Chirp)
    }
    data.Chirps[newChirpID] = Chirp{ID: int(newChirpID),
                                    Body: body,
                                    AuthorID: authorID}
    json_data,err = json.Marshal(&data)
    if err != nil{
        log.Printf("%v: %v",DataMarshalError,err)
        return Chirp{},DataMarshalError
    }
    if err := database.writeDatabase(json_data); err != nil{
        return Chirp{},err
    }
    newChirpID += 1
    return data.Chirps[newChirpID-1],err
}

func(database *db) GetChirps() ([]Chirp,error){
    json_data,err := database.loadDatabase() 
    if err != nil{
        return nil,err
    }
    data := dbStructure{}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("%v: %v",DataUnmarshalError,err)
        return nil,DataUnmarshalError
    }
    if data.Chirps == nil{
        data.Chirps = make(map[int]Chirp)
    }
    chirps := []Chirp{}
    for _,chirp := range data.Chirps{
        chirps = append(chirps,chirp)
    }
    slices.SortFunc(chirps,func (a,b Chirp) int{
                                if a.ID < b.ID{
                                    return -1
                                }else{
                                    return 1
                                }
                           })
    return chirps,nil
}

func(database *db) GetChirpByID(ID int) (Chirp,error){
    json_data,err := database.loadDatabase()
    if err != nil{
        return Chirp{},err
    }
    data := dbStructure{}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("%v: %v",DataUnmarshalError,err)
        return Chirp{},DataUnmarshalError
    }
    if data.Chirps == nil{
        data.Chirps = make(map[int]Chirp)
    }
    if chirp,ok := data.Chirps[ID]; ok{
        return chirp,nil
    }
    return Chirp{},ChirpNotFoundError
}

func(database *db) DeleteChirpByID(ID int) error{
    json_data,err := database.loadDatabase()
    if err != nil{
        return err
    }
    data := dbStructure{}
    if err := json.Unmarshal(json_data,&data); err != nil{
        log.Printf("%v: %v",DataUnmarshalError,err)
        return DataUnmarshalError
    }
    if data.Chirps == nil{
        data.Chirps = make(map[int]Chirp)
    }
    delete(data.Chirps,ID)
    json_data,err = json.Marshal(&data)
    if err != nil{
        log.Printf("%v: %v",DataMarshalError,err)
        return DataMarshalError
    }
    if err := database.writeDatabase(json_data); err != nil{
        return err
    }
    return nil
}
