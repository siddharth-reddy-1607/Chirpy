package internals

import (
	"log"
	"os"
	"slices"
	"testing"
)

var testData = []struct{body string
                        expected []chirp}{{body: "First Message",expected: []chirp{{Body: "First Message", Id : 1}}},
                                          {body: "Second Message",expected: []chirp{{Body: "First Message", Id : 1},
                                                                                    {Body: "Second Message", Id : 2}}}}
                                          
func TestCreateDB(t *testing.T){
    db,err := NewDB()
    defer db.CloseDatabase()
    if err != nil{
        t.Fatal(err)
        return
    }
    _,err = NewDB()
    if err.Error() != "There already exists a DB Connection. Be sure to close that before attempting to open a new one"{
        t.Fatalf("Not returning the correct error when DB conenction exists. Returning : %v",err)
        return
    }
    
    log.Print("Attempting to delete file")
    if err :=os.Remove(db.db_path); err != nil{
        t.Fatalf("Error while deleting database file: %v",err)
        return
    }
    log.Print("File deleted successfully")
}

func TestAddToAndGetFromDB(t *testing.T){
    db,_ := NewDB()
    defer db.CloseDatabase()
    for _,test := range testData{
        t.Log("Adding to DB")
        _,err := db.Add(test.body)
        if err != nil{
            t.Fatal(err)
            return
        }
        t.Log("Querying from DB")
        data,err := db.Query()
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
