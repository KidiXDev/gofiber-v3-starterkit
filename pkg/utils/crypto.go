package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	redigo "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type TokenClaims struct {
	jwt.RegisteredClaims
	TokenID string `json:"jti"`
}

func generateTokenID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func GenerateTokenPair(redisClient *redigo.Client, userID string) (accessToken, refreshToken string, err error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	accessTokenID := generateTokenID()
	accessClaims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenID: accessTokenID,
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = access.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	refreshTokenID := generateTokenID()
	refreshClaims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenID: refreshTokenID,
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refresh.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	ctx := context.Background()

	accessKey := "access_token:" + accessTokenID
	err = redisClient.Set(ctx, accessKey, userID, AccessTokenExpiry).Err()
	if err != nil {
		return "", "", err
	}

	refreshKey := "refresh_token:" + refreshTokenID
	err = redisClient.Set(ctx, refreshKey, userID, RefreshTokenExpiry).Err()
	if err != nil {
		return "", "", err
	}

	userRefreshKey := "user_refresh:" + userID
	redisClient.SAdd(ctx, userRefreshKey, refreshTokenID)
	redisClient.Expire(ctx, userRefreshKey, RefreshTokenExpiry)

	return accessToken, refreshToken, nil
}

func ValidateAccessToken(redisClient *redigo.Client, tokenString string) (*TokenClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	ctx := context.Background()
	accessKey := "access_token:" + claims.TokenID
	exists, err := redisClient.Exists(ctx, accessKey).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, errors.New("token has been revoked")
	}

	return claims, nil
}

func ValidateRefreshToken(redisClient *redigo.Client, tokenString string) (*TokenClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	ctx := context.Background()
	refreshKey := "refresh_token:" + claims.TokenID
	exists, err := redisClient.Exists(ctx, refreshKey).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, errors.New("refresh token has been revoked")
	}

	return claims, nil
}

func RevokeAccessToken(redisClient *redigo.Client, tokenID string) error {
	ctx := context.Background()
	accessKey := "access_token:" + tokenID
	return redisClient.Del(ctx, accessKey).Err()
}

func RevokeRefreshToken(redisClient *redigo.Client, tokenID string, userID string) error {
	ctx := context.Background()

	refreshKey := "refresh_token:" + tokenID
	err := redisClient.Del(ctx, refreshKey).Err()
	if err != nil {
		return err
	}

	userRefreshKey := "user_refresh:" + userID
	redisClient.SRem(ctx, userRefreshKey, tokenID)

	return nil
}

func RevokeAllUserTokens(redisClient *redigo.Client, userID string) error {
	ctx := context.Background()

	userRefreshKey := "user_refresh:" + userID
	refreshTokenIDs, err := redisClient.SMembers(ctx, userRefreshKey).Result()
	if err != nil {
		return err
	}

	for _, tokenID := range refreshTokenIDs {
		refreshKey := "refresh_token:" + tokenID
		redisClient.Del(ctx, refreshKey)
	}
	redisClient.Del(ctx, userRefreshKey)

	return nil
}

func RotateRefreshToken(redisClient *redigo.Client, oldClaims *TokenClaims) (accessToken, refreshToken string, err error) {
	err = RevokeRefreshToken(redisClient, oldClaims.TokenID, oldClaims.Subject)
	if err != nil {
		return "", "", err
	}

	return GenerateTokenPair(redisClient, oldClaims.Subject)
}
