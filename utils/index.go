package utils

import (
	"encoding/json"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func StructToMap(inputStruct interface{}) map[string]interface{} {

	inputJSON, err := json.Marshal(inputStruct)
	ErrorChecking(err)

	var inputMap map[string]interface{}

	err = json.Unmarshal(inputJSON, &inputMap)
	ErrorChecking(err)

	return inputMap
}

func ErrorChecking(err interface{}) {
	if err != nil {
		log.Printf("%+v", err)
	}
}

func GenerateAuthToken(email string) string {
	/*
		Method for generating the token
	*/
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	authToken, err := token.SignedString([]byte("secert"))

	ErrorChecking(err)

	return authToken
}
