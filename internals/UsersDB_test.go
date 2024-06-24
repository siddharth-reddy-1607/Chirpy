package internals

import (
	"errors"
	"log"
	"os"
	"testing"
)

func TestAddUser(t *testing.T){
    testdata := []struct{userInfo RequestUserInfo;expected ResponseUserInfo}{
                        {userInfo: RequestUserInfo{Email : "sidd@gmail.com", Password : "Sidd123!"},expected: ResponseUserInfo{ID: 1, Email: "sidd@gmail.com"}},
                        {userInfo: RequestUserInfo{Email : "niv@gmail.com", Password : "IamNivruth!"},expected: ResponseUserInfo{ID: 2, Email: "niv@gmail.com"}},
                        {userInfo: RequestUserInfo{Email : "sidd@gmail.com", Password : "DONOTCREATE!"},expected: ResponseUserInfo{}}}
    database,err := NewUsersDB()
    for _,test := range testdata{
        if err != nil{
            log.Fatalf("Error creating database connection : %v\n",err)
            return
        }
        response,err := database.AddUser(test.userInfo)
        if err != nil{
            if response != test.expected{
                log.Fatalf(`EXPECTED
                            %+v
                            FOUND
                            %+v`,test.expected,response)
                return

            }
            log.Printf("Error while add user with info %+v : %v\n",test.userInfo,err)
        }
        if response != test.expected{
            log.Fatalf(`EXPECTED
                        %+v
                        FOUND
                        %+v`,test.expected,response)
        }
    }
    log.Print("Attempting to delete file")
    if err :=os.Remove(database.db_path); err != nil{
        t.Fatalf("Error while deleting database file: %v",err)
        return
    }
    log.Print("File deleted successfully")
}

func TestLogin(t *testing.T){
    database,_:= NewUsersDB()
    addUsers := []RequestUserInfo{{Email : "sidd@gmail.com", Password : "Sidd123!"},
                                  {Email : "niv@gmail.com", Password : "IamNivruth!"}}
    for _,user := range addUsers{
        _,_ = database.AddUser(user)
    }
    type expected struct{resposne ResponseUserInfo 
                         err error}
    testdata := []struct{userInfo RequestUserInfo;Expected expected}{
                        {userInfo: RequestUserInfo{Email : "sidd@gmail.com", Password : "Sidd123!"},Expected: expected{resposne : ResponseUserInfo{ID: 1, Email: "sidd@gmail.com"}, err :nil}},
                        {userInfo: RequestUserInfo{Email : "niv@gmail.com", Password : "IamNivruth!"},Expected: expected{resposne : ResponseUserInfo{ID: 2, Email: "niv@gmail.com"}, err :nil}},
                        {userInfo: RequestUserInfo{Email : "nivnope@gmail.com", Password : "IamNivruth"},Expected: expected{resposne : ResponseUserInfo{}, err :errors.New("User Email Not Found")}},
                        {userInfo: RequestUserInfo{Email : "sidd@gmail.com", Password : "Wrong!"},Expected: expected{resposne : ResponseUserInfo{}, err :errors.New("Incorrect Password")}},
    }
    for _,test := range testdata{
        ResponseUserInfo,err := database.Login(test.userInfo)
        if ResponseUserInfo != test.Expected.resposne || (err == nil && test.Expected.err != nil) || (test.Expected.err == nil && err != nil) || (err != nil && err.Error() != test.Expected.err.Error()){
            log.Printf("Resposne : %v, Error : %v\n",ResponseUserInfo == test.Expected.resposne, err == test.Expected.err)
            log.Fatalf(`
EXPECTED
%+v with err %v
FOUND
%+v with err %v`,test.Expected.resposne,test.Expected.err,ResponseUserInfo,err)
        }
    }
}
