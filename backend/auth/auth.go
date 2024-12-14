package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"strconv"
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
	ID               int    `json:"id"`
	Username         string `json:"username"`
	Password         string `json:"-"` 
	Email            string `json:"email"`
	PasswordHash     string `json:"-"`
	TwoFactorSecret  string `json:"-"`
	TwoFactorEnabled bool   `json:"two_factor_enabled"`
}

type Session struct {
	SessionID   string    `json:"session_id"`
	AccID       int       `json:"acc_id"`
	Metadata    string    `json:"session_metadata"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiryTime  time.Time `json:"expiry_datetime"`
}

const (
	// SessionExpiry5Min for testing
	SessionExpiry5Min = 5 * time.Minute
	// SessionExpiry1Hour for production
	// SessionExpiry1Hour = 1 * time.Hour
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(user User) (string, error) {
	expiryTime := time.Now().Add(SessionExpiry5Min)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
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

func GenerateSecret() (string, error) {
	// Generate a random 20-byte secret
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(secret), nil
}

func ValidateTOTP(secret string, code string) bool {
	if len(code) != 6 {
		return false
	}
	
	inputCode, err := strconv.Atoi(code)
	if err != nil {
		return false
	}

	// Get current time window
	now := time.Now().UTC().Unix() / 30
	
	// Check current and adjacent time windows
	for delta := -1; delta <= 1; delta++ {
		if generateTOTP(secret, now+int64(delta)) == inputCode {
			return true
		}
	}
	return false
}

func generateTOTP(secret string, timeWindow int64) int {
	// Decode base32 secret
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return -1
	}

	// Create byte array of time
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(timeWindow))

	// Generate HMAC-SHA1
	h := hmac.New(sha1.New, key)
	h.Write(buf)
	sum := h.Sum(nil)

	// Get offset
	offset := sum[len(sum)-1] & 0xf

	// Generate 4-byte code
	code := binary.BigEndian.Uint32(sum[offset:offset+4])
	code &= 0x7fffffff
	code = code % 1000000

	return int(code)
}

func Enable2FA(db *sql.DB, userID int, secret string) error {
	_, err := db.Exec("UPDATE accounts SET two_factor_secret = $1, two_factor_enabled = true WHERE acc_id = $2",
		secret, userID)
	return err
}

func Disable2FA(db *sql.DB, userID int) error {
	_, err := db.Exec("UPDATE accounts SET two_factor_secret = NULL, two_factor_enabled = false WHERE acc_id = $1",
		userID)
	return err
}

func GetUser2FAStatus(db *sql.DB, userID int) (bool, string, error) {
	var enabled bool
	var secret sql.NullString
	err := db.QueryRow("SELECT two_factor_enabled, two_factor_secret FROM accounts WHERE acc_id = $1", userID).
		Scan(&enabled, &secret)
	if err != nil {
		return false, "", err
	}
	return enabled, secret.String, nil
}

func CreateSession(db *sql.DB, userID int) (*Session, error) {
	sessionID := GenerateSessionID()
	expiryTime := time.Now().Add(SessionExpiry5Min)
	
	session := &Session{
		SessionID:  sessionID,
		AccID:     userID,
		Metadata:  "{}",
		ExpiryTime: expiryTime,
	}

	_, err := db.Exec(`
		INSERT INTO sessions (session_id, acc_id, session_metadata, expiry_datetime)
		VALUES ($1, $2, $3, $4)
	`, session.SessionID, session.AccID, session.Metadata, session.ExpiryTime)

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return session, nil
}

func ValidateSession(db *sql.DB, sessionID string) (*Session, error) {
	var session Session
	err := db.QueryRow(`
		SELECT session_id, acc_id, session_metadata, created_at, expiry_datetime 
		FROM sessions 
		WHERE session_id = $1
	`, sessionID).Scan(&session.SessionID, &session.AccID, &session.Metadata, &session.CreatedAt, &session.ExpiryTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	if time.Now().After(session.ExpiryTime) {
		DeleteSession(db, sessionID)
		return nil, errors.New("session expired")
	}

	return &session, nil
}

func DeleteSession(db *sql.DB, sessionID string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE session_id = $1", sessionID)
	return err
}

func DeleteExpiredSessions(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM sessions WHERE expiry_datetime < NOW()")
	return err
}

func GenerateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
