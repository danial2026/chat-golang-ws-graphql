package decoders

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt"
)

func ParseToken(ctx context.Context, t string) (jwt.MapClaims, error) {

	// Verify
	token, parseErr := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		//HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			errMsg := fmt.Sprintf(ErrInvalidSigningMethod+"%v", token.Header["alg"])
			log.Println("ParseToken/SigningMethodHMAC/errMsg: ", errMsg)
			return nil, fmt.Errorf(errMsg)
		}
		return []byte(os.Getenv("JWTSECRET")), nil
	})

	if parseErr != nil {
		log.Println("JWT token has error while parsing: ", parseErr.Error())
		return nil, parseErr
	}

	// TODO: uncomment me after DEV is stable
	// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		/* TODO :
		[] first check client_id
		[] than check aud
		[?] check token use(token type like access or refresh)
		[] check issuer
		[] check the token expiry
		*/
		return claims, nil
	}

	return nil, errors.New(ErrInvalidToken)
}

func getUsername(claims jwt.MapClaims) string {
	if username, ok := claims["username"].(string); ok {
		return username
	}

	return ""
}

func GenerateContextAuth(claims jwt.MapClaims, token string) jwt.MapClaims {
	claims["user_id"] = claims["id"].(string)
	claims["username"] = getUsername(claims)
	claims["accessToken"] = token

	return claims
}
