package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// authRequest - структура для преобразования запроса клиента в формат JSON при аутентификации.
type authRequest struct {
	Password string `json:"password"`
}

// authResp -  структура для формирования JSON с токеном при успешной аутентификации.
type authResp struct {
	Token string `json:"token"`
}

// secret - секрет для формирования JWT токена.
var secret string

// envPassword - установленный пароль для аутентификации из переменного окружения.
var envPassword string

// SetEnvPassword - устанавливает значение envPassword, вызывается в функции main.
func SetEnvPassword() {
	envPassword = os.Getenv("TODO_PASSWORD")
	if envPassword == "" {
		log.Println("The password for authentication is not set")
		return
	}
	log.Println("The password for authentication is set")

}

// SetSecret - устанавливает значение secret, вызывается в функции main.
func SetSecret() {
	secret = os.Getenv("TODO_SECRET_KEY")
	if secret == "" {
		secret = "my_secret_key"
		log.Print("The default value of the secret key")
		//return "", fmt.Errorf("the secret key cannot be found")
	} else {
		log.Print("Secret key from a variable environment")
	}
}

// createToken - создаёт JWT токен для идентификации пользователя.
// Принимает строку - пароль, возвращает строку - JWT токен, ошибку
func createToken(password string) (string, error) {
	// Создаем хэш пароля
	hash := sha256.Sum256([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])

	// Создаём JWT
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password_hash": passwordHash, // добавляем хэш пароля в полезную нагрузку
	})

	// Получаем подписанный токен
	signedToken, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil

}

// auth - метод для идентификации пользователя при обработки запросов.
// Проверят токе из куки с валидным токеном токеном,
// если установлено значение пароля для аутентификации в переменной окружения
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля

		if len(envPassword) > 0 {
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			jwt = cookie.Value

			var valid bool
			signedToken, err := createToken(envPassword)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if signedToken == jwt {
				valid = true
			}

			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	})
}

// authHandler - хендлер для аутентификации.
// Реализует метод  Post.
// В теле запроса ожидает JSON c параметром  password.
// При соответствии пароля с парлем установленном в перемнной окружения,
// возвращает JSON с токеном для идентификации пользователя - token,
// в противном случае возвращает JSON с ошибкой - error.
func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, errorResp{Error: "Method not allowed"}, false)
		return
	}

	var buf bytes.Buffer
	var password authRequest

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &password); err != nil {
		writeJson(w, errorResp{Error: "Ошибка десериализации JSON"}, false)
		return
	}

	if envPassword == password.Password {
		taken, err := createToken(password.Password)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJson(w, authResp{Token: taken}, true)
		return
	}
	writeJson(w, errorResp{Error: "Неверный пароль"}, false)

}
