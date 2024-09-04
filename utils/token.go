package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/QingShan-Xu/xjh/cf"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct{}

func (_this *Token) GenToken(tokenStruct interface{}) (string, error) {
	data, err := json.Marshal(tokenStruct)
	if err != nil {
		return "", fmt.Errorf("TokenErr")
	}

	// 定义 JWT 的声明
	claims := jwt.MapClaims{
		"sub":  "user",                                // 主题为用户
		"exp":  time.Now().Add(72 * time.Hour).Unix(), // 过期时间为 3 天后
		"iat":  time.Now().Unix(),                     // 签发时间
		"data": data,                                  // 用户ID
	}

	// 使用声明创建新 Token
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用配置中的密钥签名 Token
	token, err := tk.SignedString([]byte(cf.TokenJWT))
	if err != nil {
		return "", fmt.Errorf("用户校验失败")
	}

	return token, nil
}

// 验证并解析 Token
func (_this *Token) ParseToken(tokenStr string) (interface{}, error) {
	jws, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(cf.TokenJWT), nil
	})

	if jws.Valid {
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return 0, fmt.Errorf("非法访问")
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return 0, fmt.Errorf("无效用户")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return 0, fmt.Errorf("登录已过期, 请重新登录")
	} else {
		return 0, err
	}

	claims, ok := jws.Claims.(jwt.MapClaims)

	if !ok {
		return 0, fmt.Errorf("无效的Token")
	}

	return claims, nil
}
