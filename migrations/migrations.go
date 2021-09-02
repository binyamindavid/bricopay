package migrations

import (
	"bricopay/helpers"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

type Account struct {
	gorm.Model
	Type    string
	Name    string
	Balance uint
	UserID  uint
}

func connectDB() *gorm.DB {

	dsn := "host=localhost port=5432 user=benjamindavid dbname=bricopay sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	helpers.HandleErr(err)

	return db
}

func createAccounts() {
	db := connectDB()

	users := []User{
		{Username: "Benjamin", Email: "binyamin.dev@gmail.com"},
		{Username: "Michael", Email: "michael.dev@gmail.com"},
	}

	for i := 0; i < len(users); i++ {
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Username))

		user := User{Username: users[i].Username, Email: users[i].Email, Password: generatedPassword}
		db.Create(&user)

		account := Account{Type: "Basic Account", Name: string(users[i].Username + "'s" + " account"), Balance: uint(10000 * int(i+1)), UserID: user.ID}
		db.Create(&account)
	}

	fmt.Println("Finished creating account")

}

func Migrate() {
	db := connectDB()
	db.AutoMigrate(&User{}, &Account{})

	createAccounts()
}
