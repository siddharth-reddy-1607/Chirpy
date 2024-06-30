package internals

import (
	"log"
	"os"
	"slices"
	"testing"
)

func TestAddToAndGetFromDB(t *testing.T){
    db,_ := NewChirpsDB()
    testData := []struct{chirpInfo RequestChirpInfo
                        expected []ResponseChirpInfo}{{chirpInfo : RequestChirpInfo{Body: "First Message",Author_Id: 1},expected : []ResponseChirpInfo{{Body: "First Message",Id: 1, AuthorId: 1}}},
                                                      {chirpInfo : RequestChirpInfo{Body: "Second Message",Author_Id: 1},expected : []ResponseChirpInfo{{Body: "First Message",Id: 1, AuthorId: 1},{Body: "Second Message",Id: 2, AuthorId: 1}}}}
                                        
    for _,test := range testData{
        t.Log("Adding to DB")
        _,err := db.AddChirp(test.chirpInfo)
        if err != nil{
            t.Fatal(err)
            return
        }
        t.Log("Querying from DB")
        data,err := db.QueryChirps()
        if err != nil{
            t.Fatal(err)
            return
        }
        if slices.Equal(data, test.expected) == false{
            t.Fatalf("\nEXPECTED \n %v \n FOUND \n %v \n",test.expected,data)
            return
        }
    }
    log.Print("Attempting to delete file")
    if err :=os.Remove(db.db_path); err != nil{
        t.Fatalf("Error while deleting database file: %v",err)
        return
    }
    log.Print("File deleted successfully")
}

func TestQueryChirpByID(t *testing.T){
    db,_ := NewChirpsDB()
    dataToAdd := []RequestChirpInfo{{Body: "First Message",Author_Id: 1},{Body: "Second Message",Author_Id: 1}}
    testData := []struct{ID int
    expected ResponseChirpInfo}{{ID : 1, expected: ResponseChirpInfo{Body : "First Message", AuthorId: 1, Id: 1}},
                    {ID : 2, expected: ResponseChirpInfo{Body : "Second Message", AuthorId: 1, Id: 2}}}
    for _,test := range dataToAdd{
        db.AddChirp(test)
    }
    for _,test := range testData{
        c,err := db.QueryChirpByID(test.ID)
        if err != nil{
            log.Fatalf("Error getting chirp by ID %d : %v\n",test.ID,err)
            return
        }
        if c != test.expected{
            log.Fatalf(`\nEXPECTED
            %+v
            FOUND,
            %+v`,test.expected,c)
            return
        }
    }
    log.Print("Attempting to delete file")
    if err :=os.Remove(db.db_path); err != nil{
        t.Fatalf("Error while deleting database file: %v",err)
        return
    }
    log.Print("File deleted successfully")
}
