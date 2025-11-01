// Package auth provides functions related to password hashing and session
// tokens.
package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/johndosdos/chatter/internal/database"
)

// The ContextKey type is meant for passing userID as key for
// context.WithValue.
type ContextKey string

// UserIDKey implements the ContextKey type.
const UserIDKey ContextKey = "userId"

// HashPassword returns the hashed password created using the argon2id
// package.
func HashPassword(password string) (string, error) {
	hashedPw, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("internal/auth: pw hash failed: %w", err)
	}

	return hashedPw, nil
}

// CheckPasswordHash compares the password and the hash. It returns true
// when they match, otherwise it returns false.
func CheckPasswordHash(password, hash string) (bool, error) {
	isMatch, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("internal/auth: pw and hash comparison failed: %w", err)
	}
	if !isMatch {
		return false, errors.New("internal/auth: pw and hash do not match")
	}

	return isMatch, nil
}

// MakeJWT returns a JSON Web Token string to be used as an acess token
// for client session.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    os.Getenv("JWT_ISS"),
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
	})

	return token.SignedString([]byte(tokenSecret))
}

// ValidateJWT tries to validate the access token. It returns the user id
// as a uuid.UUID type. The returned error is a uuid.Parse error.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(*jwt.Token) (any, error) { return []byte(tokenSecret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("internal/auth: failed to parse token: %w", err)
	}

	if !token.Valid {
		return uuid.UUID{}, errors.New("internal/auth: token is invalid")
	}

	if claims.Subject == "" {
		return uuid.UUID{}, errors.New("subject claim is missing")
	}

	userid, _ := token.Claims.GetSubject()
	return uuid.Parse(userid)
}

// MakeRefreshToken returns a refresh token string, while also storing the
// token to the database.
func MakeRefreshToken(ctx context.Context, db *database.Queries) (string, error) {
	rnd := make([]byte, 32)

	// rand.Read() never returns an error.
	_, _ = rand.Read(rnd)
	rndStr := hex.EncodeToString(rnd)

	userID := ctx.Value(UserIDKey).(uuid.UUID)
	refreshTokenExp := time.Now().UTC().AddDate(0, 0, 7)
	refreshToken, err := db.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:     rndStr,
		CreatedAt: pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		ExpiresAt: pgtype.Timestamptz{Time: refreshTokenExp, Valid: true},
	})
	if err != nil {
		return "", fmt.Errorf("database error: %w", err)
	}

	return refreshToken.Token, nil
}
