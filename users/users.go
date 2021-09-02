package users

import (
	"bricopay/helpers"
	"bricopay/interfaces"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func prepareToken(user *interfaces.User) string {
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute * 60).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	return token
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount) map[string]interface{} {

	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}

	var token = prepareToken(user)
	var response = map[string]interface{}{"message": "all is ok"}
	response["jwt"] = token
	response["data"] = responseUser

	return response

}

func Login(username string, password string) map[string]interface{} {
	// Add validation to login
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: password, Valid: "password"},
		})
	if valid {
		// Connect DB
		db := helpers.ConnectDB()
		user := &interfaces.User{}
		if err := db.Where("username = ? ", username).First(&user).Error; err != nil {
			return map[string]interface{}{"message": "User not found"}
		}
		// Verify password
		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "Wrong password"}
		}
		// Find  user accounts
		accounts := []interfaces.ResponseAccount{}
		db.Table("accounts").Select("id, name, balance").Where("user_id = ? ", user.ID).Scan(&accounts)

		var response = prepareResponse(user, accounts)

		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}
}

func Register(username string, email string, password string) map[string]interface{} {
	// Add validation to registration
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: email, Valid: "email"},
			{Value: password, Valid: "password"},
		})
	if valid {
		// Create registration logic
		// Connect DB
		db := helpers.ConnectDB()
		generatedPassword := helpers.HashAndSalt([]byte(password))
		user := &interfaces.User{Username: username, Email: email, Password: generatedPassword}
		db.Create(&user)

		account := &interfaces.Account{Type: "Basic Account", Name: string(username + "'s" + " account"), Balance: 0, UserID: user.ID}
		db.Create(&account)

		accounts := []interfaces.ResponseAccount{}
		respAccount := interfaces.ResponseAccount{ID: account.ID, Name: account.Name, Balance: account.Balance}
		accounts = append(accounts, respAccount)
		var response = prepareResponse(user, accounts)

		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}

}
