package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)


type JWTToken struct {
	config *Config
}


type jwtClaim struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
	Exp int64 `json:"exp"`
}

func NewJWTToken(config *Config) *JWTToken {
	return &JWTToken{ config: config }
}

func (j *JWTToken) CreateToken(userID int64) (string, error) {
	claims := jwtClaim{
		UserID: userID,
		Exp: time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims);

	tokenString, err := token.SignedString([]byte(j.config.SigningKey));
	if err != nil {
		return "", err
	}

	return string(tokenString), nil
}

func (j *JWTToken) VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _,ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid auth token provided");
		}

		return []byte(j.config.SigningKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid auth token provided");
	}

	claim, ok := token.Claims.(*jwtClaim);

	if !ok {
		return 0, fmt.Errorf("invalid auth token provided");
	}
	
	if claim.Exp < time.Now().Unix() {
		return 0, fmt.Errorf("provided expired token");
	}


	return claim.UserID, nil
}

