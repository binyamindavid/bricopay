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

	db := helpers.ConnectDB()
	user := &interfaces.User{}

	if err := db.Where("username = ? ", username).First(&user).Error; err != nil {
		return map[string]interface{}{"message": "User not found"}
	}

	// Verify Password
	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if passwordErr == bcrypt.ErrMismatchedHashAndPassword && passwordErr != nil {
		return map[string]interface{}{"message": "Wrong Password"}
	}

	accounts := []interfaces.ResponseAccount{}
	db.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

	var response = prepareResponse(user, accounts)

	return response

}
