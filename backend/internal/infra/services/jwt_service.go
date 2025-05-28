package services

import (
	"errors"
	"fmt"
	"main/config"
	domain "main/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
	config *config.Config `di.inject:"config"`
}

func (s *JWTService) GenerateToken(account *domain.Account) (string, error) {
	tokenLifespan := s.config.JwtTokenLifespanHour
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = account.Id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Settings.ApiSecret))
}

func (s *JWTService) ValidateToken(token string) error {
	parsedToken, err := s.parseToken(token)

	if err != nil {
		return err
	}

	_, ok := parsedToken.Claims.(jwt.MapClaims)
	if ok && parsedToken.Valid {
		return nil
	}

	return errors.New("Invalid token provided")
}

func (s *JWTService) GetAccountId(token string) (domain.AccountId, error) {
	parsedToken, err := s.parseToken(token)

	if err != nil {
		return domain.AccountId(0), err
	}

	claims, _ := parsedToken.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	return domain.AccountId(userId), nil
}

func (s *JWTService) parseToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Settings.ApiSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error parsing token: %w", err)
	}

	return parsedToken, nil
}
