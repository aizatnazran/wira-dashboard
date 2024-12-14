package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey []byte

func InitJWTKey(secret string) {
	jwtKey = []byte(secret)
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"-"` 
	Email        string `json:"email"`
	PasswordHash string `json:"-"` 
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(user User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func CreateUser(db *sql.DB, username, password, email string) error {
	// Check if username already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking username: %v", err)
	}
	if exists {
		return errors.New("username already exists")
	}

	// Check if email already exists
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking email: %v", err)
	}
	if exists {
		return errors.New("email already exists")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	query := `
		INSERT INTO accounts (username, password_hash, email)
		VALUES ($1, $2, $3)
	`
	_, err = db.Exec(query, username, hashedPassword, email)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

func AuthenticateUser(db *sql.DB, username, password string) (*User, error) {
	user := &User{}
	query := `
		SELECT acc_id, username, password_hash, email
		FROM accounts
		WHERE username = $1
	`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid username or password")
		}
		return nil, fmt.Errorf("error querying user: %v", err)
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}
