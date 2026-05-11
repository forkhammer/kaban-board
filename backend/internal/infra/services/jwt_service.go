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
	return token.SignedString([]byte(s.config.ApiSecret + account.JwtSalt))
}

func (s *JWTService) ValidateToken(token string, salt string) error {
	parsedToken, err := s.parseToken(token, salt)
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
	parser := jwt.NewParser()
	parsedToken, _, err := parser.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return domain.AccountId(0), fmt.Errorf("Error parsing token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return domain.AccountId(0), errors.New("Invalid token claims")
	}

	v, ok := claims["id"].(float64)
	if !ok {
		return domain.AccountId(0), fmt.Errorf("id claim is not a number, got %T", claims["id"])
	}
	return domain.AccountId(uint(v)), nil
}

func (s *JWTService) parseToken(token, salt string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.config.ApiSecret + salt), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error parsing token: %w", err)
	}

	return parsedToken, nil
}
