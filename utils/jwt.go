package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"oe02_go_tam/constant"
	"os"
	"time"
)

var jwtSecret []byte

func InitJWTSecret() error {
	if err := godotenv.Load(); err != nil {
		return errors.New(constant.T(constant.ErrLoadEnvFailed.Error()))
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return errors.New(constant.T(constant.ErrJWTSecretNotSet.Error()))
	}

	jwtSecret = []byte(secret)
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (uint, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(constant.T(constant.UnexpectedSigningMethod.Error()))
		}
		return jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, "", errors.New(constant.T(constant.TokenExpired.Error()))
		}
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", errors.New(constant.T(constant.InvalidToken.Error()))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New(constant.T(constant.InvalidTokenClaim.Error()))
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", errors.New(constant.T(constant.UserIdNotFoundInToken.Error()))
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", errors.New(constant.T(constant.RoleNotFoundInToken.Error()))
	}

	return uint(userIDFloat), role, nil
}
