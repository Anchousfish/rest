package service

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/Anchousfish/rest/models"
	"github.com/Anchousfish/rest/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
)

const (
	salt       = "jehrioiefmdv,nsfgjr"
	tokenTTL   = 12 * time.Hour
	signingKey = "lakjrotiuotijsvmsckdfkljgqhwgjokgklflk"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func newAuthService(rp repository.Authorization) *AuthService {
	return &AuthService{repo: rp}
}

func (a *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return a.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (a *AuthService) GenerateToken(username, password string) (string, error) {

	user, err := a.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", errors.Wrap(err, "user not found")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: user.Id,
	})
	return token.SignedString([]byte(signingKey))
}

func (a AuthService) ParseToken(accesstoken string) (userid int, err error) {
	token, err := jwt.ParseWithClaims(accesstoken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of the type *tokenClaims")
	}

	return claims.UserId, nil

}
