package account

import (
	"errors"
	"fmt"
	"html"
	"main/config"
	"main/repository"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	accountRespository repository.AccountRepositoryInterface `di.inject:"accountRepository"`
}

func (s *AccountService) GetPasswordHash(rawPass string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPass), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return html.EscapeString(strings.TrimSpace(string(hashedPassword))), nil
}

func (s *AccountService) RegisterAccount(username, password string) (*Account, error) {
	err := s.accountRespository.GetAccountByUsername(&Account{}, username)

	if err == nil {
		return nil, errors.New("Такой пользователь уже существует")
	}

	passwordHash, err := s.GetPasswordHash(password)

	if err != nil {
		return &Account{}, err
	}

	account := &Account{
		Name:     username,
		Username: username,
		Password: passwordHash,
		IsActive: false,
	}
	err = s.accountRespository.CreateAccount(account)

	return account, err
}

func (s *AccountService) GetTokenByCredentials(username, password string) (string, error) {
	var account Account
	err := s.accountRespository.GetAccountByUsername(&account, username)

	if err != nil {
		return "", errors.New("Такой пользователь не найден")
	}

	if !account.IsActive {
		return "", errors.New("Пользователь не активирован")
	}

	err = s.VerifyPassword(password, account.Password)

	if err != nil {
		return "", errors.New("Неверный пароль")
	}

	token, err := s.GenerateToken(&account)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AccountService) VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *AccountService) GenerateToken(account *Account) (string, error) {
	tokenLifespan := config.Settings.JwtTokenLifespanHour
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = account.Id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Settings.ApiSecret))
}

func (s *AccountService) ValidateToken(c *gin.Context) error {
	token, err := s.GetToken(c)

	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}

	return errors.New("Invalid token provided")
}

func (s *AccountService) GetToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := s.getTokenFromRequest(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Settings.ApiSecret), nil
	})
	return token, err
}

func (s *AccountService) getTokenFromRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func (s *AccountService) GetCurrentAccountFromContext(c *gin.Context) (*Account, error) {
	err := s.ValidateToken(c)
	if err != nil {
		return nil, err
	}
	token, _ := s.GetToken(c)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	var user Account
	err = s.accountRespository.GetAccountById(&user, userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
