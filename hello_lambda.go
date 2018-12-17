package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tidwall/gjson"

        "time"
        "github.com/dgrijalva/jwt-go"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, evt json.RawMessage) (events.APIGatewayProxyResponse, error) {

	signingKey := os.Getenv("SECRET")

	passedToken := gjson.Get(string(evt), "jwt").String()
	secret := gjson.Get(string(evt), "token")

	returnVal := events.APIGatewayProxyResponse{
		Body:       "",
		StatusCode: 200,
	}

	if secret.String() != os.Getenv("SECRET") {
            createdToken, err := ExampleNew([]byte(signingKey))
            if err == nil {
	        returnVal.Body = "bad token, here is a working jwt: " + createdToken
	        returnVal.StatusCode = 401
	    } else {
	        returnVal.Body = "error creating token"
	        returnVal.StatusCode = 500
            }
	} else {
	    fmt.Println("Hello Lambda")
            token, err := jwt.Parse(passedToken, func(token *jwt.Token) (interface{}, error) {
                return []byte(signingKey), nil
            })
            if err == nil && token.Valid {
                claims := token.Claims.(jwt.MapClaims)
                // TODO test for presence of claim
                returnVal.Body = claims["foo"].(string)
            } else {
                returnVal.Body = "jwt not valid or missing required claim"
	        returnVal.StatusCode = 400
            }
	}
	return returnVal, nil
}

func ExampleNew(mySigningKey []byte) (string, error) {
    // Create the token
    token := jwt.New(jwt.SigningMethodHS256)
    // Set some claims
    claims := make(jwt.MapClaims)
    claims["foo"] = "bar"
    claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
    token.Claims = claims
    // Sign and get the complete encoded token as a string
    tokenString, err := token.SignedString(mySigningKey)
    return tokenString, err
}

