package internals

import (
	"log"
	"testing"
    "slices"
)

func TestAddGetChirps(t *testing.T){
    db,err := NewDB("database.json")
    if err != nil{
        t.Fatalf("Error while creating DB: %v",err)
    }
    type parameters struct{
        body string
        authorID int
        expected Chirp
    }
    testData := []parameters{}
    getChirpsCallExpected := []Chirp{}
    testData = append(testData,parameters{body : "First Message",authorID : 1, expected : Chirp{ID : 1,
                                                                                                Body :"First Message",
                                                                                                AuthorID : 1}})
    testData = append(testData,parameters{body : "Second Message",authorID : 1, expected : Chirp{ID : 2,
                                                                                                Body :"Second Message",
                                                                                                AuthorID : 1}})
    testData = append(testData,parameters{body : "Third Message",authorID : 2, expected : Chirp{ID : 3,
                                                                                                Body :"Third Message",
                                                                                                AuthorID : 2}})
    for _,data := range testData{
        chirp,err := db.AddChirp(data.body,data.authorID)
        if err != nil{
            log.Fatalf("Error while adding chirp: %v",err)
            if err := db.EraseDB(); err != nil{
                t.Logf("Error while erasing DB: %v\n",err)
            }
        }
        if chirp != data.expected{
            t.Logf("Validation mismatch during AddChirps() ")
            t.Fatalf(
            `EXPECTED
             %+v
             FOUND
             %+v`,data.expected,chirp)
             if err := db.EraseDB(); err != nil{
                 t.Logf("Error while erasing DB: %v\n",err)
             }
        }
        getChirpsCallExpected = append(getChirpsCallExpected, chirp)
        chirps,err := db.GetChirps()
        if err != nil{
            log.Fatalf("Error while getting chirps: %v",err)
        }
        if slices.Equal(getChirpsCallExpected,chirps) == false{
            t.Logf("Validation mismatch during GetChirps() ")
            t.Fatalf(
            `EXPECTED
             %+v
             FOUND
             %+v`,getChirpsCallExpected,chirps)
        }
    }
}

func TestGetChirpByID(t *testing.T){
    db,err := NewDB("database.json")
    if err != nil{
        t.Fatalf("Error while creating DB: %v",err)
    }
    type parameters struct{
        ID int
        expected Chirp
    }
    testData := []parameters{}
    testData = append(testData,parameters{ID : 1, expected : Chirp{ID : 1,
                                                                  Body :"First Message",
                                                                  AuthorID : 1}})
    testData = append(testData,parameters{ID : 3, expected : Chirp{ID : 3,
                                                                  Body :"Third Message",
                                                                  AuthorID : 2}})
    for _,data := range testData{
        chirp,err := db.GetChirpByID(data.ID) 
        if err != nil{
            log.Fatalf("Error while getting chirp with ID %d: %v",err,data.ID)
            if err := db.EraseDB(); err != nil{
                t.Logf("Error while erasing DB: %v\n",err)
            }
        }
        if chirp != data.expected{
            t.Logf("Validation mismatch during GetChirpByID() ")
            t.Fatalf(
            `EXPECTED
             %+v
             FOUND
             %+v`,data.expected,chirp)
             if err := db.EraseDB(); err != nil{
                 t.Logf("Error while erasing DB: %v\n",err)
             }
        }
    }
}

func TestDeleteChirpByID(t *testing.T){
    db,err := NewDB("database.json")
    if err != nil{
        t.Fatalf("Error while creating DB: %v",err)
    }
    type parameters struct{
        ID int
        expected []Chirp
    }
    testData := []parameters{}
    testData = append(testData,parameters{ID : 1, expected : []Chirp{Chirp{ID : 2,
                                                                           Body :"Second Message",
                                                                           AuthorID : 1},
                                                                     Chirp{ID : 3,
                                                                           Body :"Third Message",
                                                                           AuthorID : 2}}})
    testData = append(testData,parameters{ID : 3, expected : []Chirp{Chirp{ID : 2,
                                                                     Body :"Second Message",
                                                                     AuthorID : 1}}})
    for _,data := range testData{
        err := db.DeleteChirpByID(data.ID) 
        if err != nil{
            log.Fatalf("Error while getting chirp with ID %d: %v",err,data.ID)
            if err := db.EraseDB(); err != nil{
                t.Logf("Error while erasing DB: %v\n",err)
            }
        }
        chirps,err := db.GetChirps()
        if slices.Equal(chirps,data.expected) == false{
            t.Logf("Validation mismatch during DeleteChirpByID() ")
            t.Fatalf(
            `EXPECTED
             %+v
             FOUND
             %+v`,data.expected,chirps)
             if err := db.EraseDB(); err != nil{
                 t.Logf("Error while erasing DB: %v\n",err)
             }
        }
    }
}




