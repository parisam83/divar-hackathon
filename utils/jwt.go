package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/configs"
	"github.com/golang-jwt/jwt"
)

type JWTManager struct {
	secretKey []byte
}

func NewJWTManager(config configs.JWTConfig) *JWTManager {
	return &JWTManager{
		secretKey: []byte(config.JwtSecret),
	}
}
func (j *JWTManager) CreateJwtToken(userId string) string {
	claims := jwt.MapClaims{
		"user_id": userId,
		// "user_status": userStatus,
		// "post_token":  postId,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token, err := jwtToken.SignedString(j.secretKey)
	if err != nil {
		log.Println("could not create jwt token")
	}
	return token

}

func (j *JWTManager) JWTMiddlewear(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Parse token
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return j.secretKey, nil
		})

		// Handle parsing errors
		if err != nil {
			if err.Error() == "Token is expired" {
				http.Error(w, "Unauthorized: Token expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			}
			return
		}

		// Validate token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Add claims to request context
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			// ctx = context.WithValue(ctx, "post_token", claims["post_token"])
			// ctx = context.WithValue(ctx, "user_status", claims["user_status"])
			next(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		}

	}

}
