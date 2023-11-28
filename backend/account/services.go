package account

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"html"
	"main/config"
	"main/db"
	"strings"
	"time"
)

func GetPasswordHash(rawPass string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPass), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return html.EscapeString(strings.TrimSpace(string(hashedPassword))), nil
}

func RegisterAccount(username, password string) (*Account, error) {
	_, err := GetAccountByUsername(username)

	if err == nil {
		return nil, errors.New("Такой пользователь уже существует")
	}

	passwordHash, err := GetPasswordHash(password)

	if err != nil {
		return &Account{}, err
	}

	account := &Account{
		Name:     username,
		Username: username,
		Password: passwordHash,
		IsActive: false,
	}
	result := db.DefaultConnection.Db.Create(account)

	return account, result.Error
}

func GetTokenByCredentials(username, password string) (string, error) {
	account, err := GetAccountByUsername(username)

	if err != nil {
		return "", errors.New("Такой пользователь не найден")
	}

	if !account.IsActive {
		return "", errors.New("Пользователь не активирован")
	}

	err = VerifyPassword(password, account.Password)

	if err != nil {
		return "", errors.New("Неверный пароль")
	}

	token, err := GenerateToken(account)

	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateToken(account *Account) (string, error) {
	tokenLifespan := config.Settings.JwtTokenLifespanHour
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = account.Id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Settings.ApiSecret))
}

func GetAccountByUsername(username string) (*Account, error) {
	var account Account

	result := db.DefaultConnection.Db.Where(Account{Username: username}).First(&account)

	return &account, result.Error
}

func GetAccountById(id uint) (*Account, error) {
	var account Account

	result := db.DefaultConnection.Db.Where(Account{Id: id}).First(&account)

	if result.Error != nil {
		return nil, result.Error
	}

	return &account, nil
}

func ValidateToken(c *gin.Context) error {
	token, err := GetToken(c)

	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}

	return errors.New("Invalid token provided")
}

func GetToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Settings.ApiSecret), nil
	})
	return token, err
}

func getTokenFromRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func GetCurrentAccountFromContext(c *gin.Context) (*Account, error) {
	err := ValidateToken(c)
	if err != nil {
		return nil, err
	}
	token, _ := GetToken(c)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	user, err := GetAccountById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
