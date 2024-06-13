package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"gorm.io/gorm"
)

type UserClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for authenticated users
func GenerateToken(u *model.User) (string, error) {
	// Convert roles to a slice of role names for inclusion in the token.

	claims := &UserClaims{
		ID:       u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // Set the expiration time
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("your_secret_key")) // Use the same key as in the middleware
	if err != nil {
		fmt.Println("Error generating token: ", err)
		return "", err
	}

	return t, nil
}


// Tạo Refresh Token
func GenerateRefreshToken(user *model.User) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["userID"] = user.ID
    claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix() // Refresh token có thời hạn là 30 ngày

    refreshToken, err := token.SignedString([]byte("refresh_token_secret")) // Đổi "refresh_token_secret" thành secret của bạn
    if err != nil {
        return "", err
    }
    return refreshToken, nil
}

// Xác Thực Refresh Token
func ValidateRefreshToken(tokenString string) (*model.User, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Kiểm tra loại token và trả về secret key
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte("refresh_token_secret"), nil // Đổi "refresh_token_secret" thành secret của bạn
    })
    if err != nil {
        return nil, err
    }
    // Kiểm tra token và lấy thông tin user từ claims
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := claims["userID"].(uint) // Giả sử userID được lưu trong claims
        user := &model.User{
			Model: gorm.Model{
				ID: userID,
			},
		}   
        return user, nil
    }
    return nil, fmt.Errorf("invalid refresh token")
}