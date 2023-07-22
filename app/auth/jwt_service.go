package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidJWTSigningMethod = errors.New("Invalid JWT signing method")
	ErrInvalidJWTClaims        = errors.New("Invalid JWT claims")
)

type accessTokenClaims struct {
	jwt.RegisteredClaims
	Scopes string `json:"scopes"`
}

func validateUserAccessToken(ctx context.Context, tokenStr string) (user User, err error) {
	claims := &accessTokenClaims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			_, isAccSigningMethod := t.Method.(*jwt.SigningMethodHMAC)
			if !isAccSigningMethod {
				return nil, ErrInvalidJWTSigningMethod
			}
			return jwtAccessTokenSecret, nil
		},
		jwt.WithIssuer(jwtIssuer),
		jwt.WithAudience(jwtIssuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return
	}

	claims, isValidClaims := token.Claims.(*accessTokenClaims)
	if !isValidClaims || !strings.Contains(claims.Scopes, "ROLE_USER") {
		return user, ErrInvalidJWTClaims
	}

	userId, err := validateUserId(claims.Subject)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to validate user access token")
		return
	}

	defer tx.Rollback(ctx)

	user, err = findUserById(ctx, tx, userId)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to validate user access token")
		return
	}

	return user, nil
}

func generateUserAccessToken(userId string) (token string, err error) {
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Audience:  []string{jwtIssuer},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			ID:        ulid.Make().String(),
			Subject:   userId,
		},
		Scopes: "ROLE_USER",
	})

	token, err = tokenObj.SignedString(jwtAccessTokenSecret)
	if err != nil {
		log.Err(err).Msg("Failed to generate user access token")
		return
	}

	return token, nil
}
