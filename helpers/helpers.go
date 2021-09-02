package helpers

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	HandleErr(err)

	return string(hashed)
}

func ConnectDB() *gorm.DB {

	dsn := "host=localhost port=5432 user=benjamindavid dbname=bricopay sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	HandleErr(err)

	return db
}
