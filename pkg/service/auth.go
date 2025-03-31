package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ponomare0v/todo-go-app/pkg/models"
	"github.com/ponomare0v/todo-go-app/pkg/repository"
)

const (
	salt       = "fahfashj23h4hvxgdau7434"
	tokenTTL   = 12 * time.Hour
	signingKey = "effcjafsc5638c2xdw82323xfkwiwr34u5b3i"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

//
//

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

//
//

func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

//
//

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

//
//

// метод структуры AuthService, где вызовем функцию из библиотеки jwt ParseWithClaims,
// которая принимает сам токен, структуру Claims и функцию, которая возвращает ключ подпись или ошибку.
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		//В этой функции нужно проверить метод подписи токена, если это не HMAC, то возвращем ошибку,
		//  а если все окей то возвращаем ключ-подпись (signingKey).
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	//Функция parseWithClaims возвращет объект токена, в котором есть поле Claims типа интерфейс,
	//  приведим его к собственной структуре и проверяем всё ли хорошо.
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	//Если все же успешно распарсили токен - вернем значение id пользователя.
	return claims.UserId, nil
}

//
//

// функцию хэширования пароля generatePasswordHash используя алгоритм sha1, также - добавляем к паролю некий набор случайных символов, перед тем как
// проводить процедуру хэширования. Этот набор символов еще называют солью - конст salt со случайными символами.
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
