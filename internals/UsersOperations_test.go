package internals

import (
	"testing"
    "log"
    "slices"
)

func TestAddGetUsers(t *testing.T){
    db,err := NewDB("database.json")
    if err != nil{
        t.Fatalf("Error while creating DB: %v",err)
    }
    type parameters struct{email string
                           encryptedPassword string
                           expected User}
    testData := []parameters{}
    testData = append(testData,parameters{"siddharth@gmail.com","Sid@18",User{ID : 1,
                                                                              Email : "siddharth@gmail.com",
                                                                              EncryptedPassword: "Sid@18"}})
    testData = append(testData,parameters{"nivruth@gmail.com","Niv@2004",User{ID : 2,
                                                                              Email : "nivruth@gmail.com",
                                                                              EncryptedPassword: "Niv@2004"}})
    testData = append(testData,parameters{"sandy@gmail.com","SandyWorksHard@891979",User{ID : 3,
                                                                                         Email : "sandy@gmail.com",
                                                                                         EncryptedPassword: "SandyWorksHard@891979"}})
    getUsersCallExpected := []User{}
    for _,data := range testData{
        user,err := db.AddUser(data.email,data.encryptedPassword)
        if err != nil{
            log.Fatalf("Error while adding chirp: %v",err)
            if err := db.EraseDB(); err != nil{
                t.Logf("Error while erasing DB: %v\n",err)
            }
        }
        if user != data.expected{
            t.Fatalf(`
            EXPECTED
            %v
            FOUND
            %v
            `,data.expected,user)
        }
        getUsersCallExpected = append(getUsersCallExpected,user)
        users,err := db.GetUsers()
        if err != nil{
            log.Fatalf("Error while adding chirp: %v",err)
            if err := db.EraseDB(); err != nil{
                t.Logf("Error while erasing DB: %v\n",err)
            }
        }
        if slices.Equal(users,getUsersCallExpected) == false{
            t.Fatalf(`
            EXPECTED
            %v
            FOUND
            %v
            `,getUsersCallExpected,users)
        }
    }
}

func TestGetUserByID(t *testing.T){
    db,err := NewDB("database.json")
    if err != nil{
        t.Fatalf("Error while creating DB: %v",err)
    }
    type expected struct{
        user User
        err error
    }
    type parameters struct{
        ID int
        Expected expected
    }
    testData := []parameters{}
    testData = append(testData, parameters{1,expected{User{Email:"siddharth@gmail.com",
                                                           EncryptedPassword:"Sid@18",
                                                           ID: 1},
                                                           nil}})
    testData = append(testData, parameters{3,expected{User{Email:"sandy@gmail.com",
                                                           EncryptedPassword:"SandyWorksHard@891979",
                                                           ID: 3},
                                                           nil}})
    testData = append(testData, parameters{4,expected{User{},
                                                           UserNotFoundError}})

    for _,data := range testData{
        user,err := db.GetUserByID(data.ID)
        if err != data.Expected.err || user != data.Expected.user{
            t.Fatalf(`
            EXPECTED
            %+v
            FOUND
            %+v
            `,data.Expected,struct{user User;err error}{user: user, err: err})
        }
    }
}

func TestGetUserByEmail(t *testing.T){
    db,err := NewDB("database.json")
    if err != nil{
        t.Fatalf("Error while creating DB: %v",err)
    }
    type expected struct{
        user User
        err error
    }
    type parameters struct{
        email string
        Expected expected
    }
    testData := []parameters{}
    testData = append(testData, parameters{"siddharth@gmail.com",expected{User{Email:"siddharth@gmail.com",
                                                                          EncryptedPassword:"Sid@18",
                                                                          ID: 1},
                                                                          nil}})
    testData = append(testData, parameters{"sandy@gmail.com",expected{User{Email:"sandy@gmail.com",
                                                                           EncryptedPassword:"SandyWorksHard@891979",
                                                                           ID: 3},
                                                                           nil}})
    testData = append(testData, parameters{"whoknows@gmail.com",expected{User{},
                                                                         UserNotFoundError}})

    for _,data := range testData{
        user,err := db.GetUserByEmail(data.email)
        if err != data.Expected.err || user != data.Expected.user{
            t.Fatalf(`
            EXPECTED
            %+v
            FOUND
            %+v
            `,data.Expected,struct{user User;err error}{user: user, err: err})
        }
    }
}

func TestUpdateUser(t *testing.T){
    db,err := NewDB("database.json")
    if err != nil{
        t.Fatalf("Error while creating DB: %v",err)
    }
    type expected struct{
        user User
        err error
    }
    type parameters struct{
        ID int
        email string
        encryptedPassword string
        Expected expected
    }
    testData := []parameters{}
    testData = append(testData, parameters{1,"sidnew@gmail.com","Sid@18",expected{User{Email:"sidnew@gmail.com",
                                                                                       EncryptedPassword:"Sid@18",
                                                                                       ID: 1},
                                                                                  nil}})
    testData = append(testData, parameters{3,"budiki@gmail.com","Budiki",expected{User{Email:"budiki@gmail.com",
                                                                                      EncryptedPassword:"Budiki",
                                                                                      ID: 3},
                                                                                  nil}})
    testData = append(testData, parameters{4,"sidnew@gmail.com","Sid@18",expected{User{},
                                                                                  UserNotFoundError}})
    for _,data := range testData{
        user,err := db.UpdateUser(data.ID,data.email,data.encryptedPassword)
        if err != data.Expected.err || user != data.Expected.user{
            t.Fatalf(`
            EXPECTED
            %+v
            FOUND
            %+v
            `,data.Expected,struct{user User;err error}{user: user, err: err})
        }
    }
}

func TestUpdgradeUser(t *testing.T){
    db,err := NewDB("database.json")
    if err != nil{
        t.Fatalf("Error while creating DB: %v",err)
    }
    type expected struct{
        user User
        err error
    }
    type parameters struct{
        ID int
        Expected expected
    }
    testData := []parameters{}
    testData = append(testData, parameters{1,expected{User{Email:"sidnew@gmail.com",
                                                           EncryptedPassword:"Sid@18",
                                                           IsChirpyRed: true,
                                                           ID: 1},
                                                      nil}})
    testData = append(testData, parameters{3,expected{User{Email:"budiki@gmail.com",
                                                           EncryptedPassword:"Budiki",
                                                           IsChirpyRed: true,
                                                           ID: 3},
                                                      nil}})
    testData = append(testData, parameters{4,expected{User{},
                                                      UserNotFoundError}})
    for _,data := range testData{
        user,err := db.UpgradeUser(data.ID)
        if err != data.Expected.err || user != data.Expected.user{
            t.Fatalf(`
            EXPECTED
            %+v
            FOUND
            %+v
            `,data.Expected,struct{user User;err error}{user: user, err: err})
        }
    }
}
