package migrations

import (
	"bricopay/helpers"
	"bricopay/interfaces"
	"fmt"
)

func createAccounts() {
	db := helpers.ConnectDB()

	users := &[2]interfaces.User{
		{Username: "Benjamin", Email: "binyamin.dev@gmail.com"},
		{Username: "Michael", Email: "michael.dev@gmail.com"},
	}

	for i := 0; i < len(users); i++ {
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Username))

		user := &interfaces.User{Username: users[i].Username, Email: users[i].Email, Password: generatedPassword}
		db.Create(&user)

		account := &interfaces.Account{Type: "Basic Account", Name: string(users[i].Username + "'s" + " account"), Balance: uint(10000 * int(i+1)), UserID: user.ID}
		db.Create(&account)
	}

	fmt.Println("Finished creating account")

}

func Migrate() {
	User := &interfaces.User{}
	Account := &interfaces.Account{}

	db := helpers.ConnectDB()
	db.AutoMigrate(User, Account)

	createAccounts()
}
