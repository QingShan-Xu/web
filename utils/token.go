package utils

// type Token struct{}

// func (_this *Token) GenToken(id int) (string, error) {

// 	// 定义 JWT 的声明
// 	claims := jwt.MapClaims{
// 		"sub":     "user",                                // 主题为用户
// 		"exp":     time.Now().Add(72 * time.Hour).Unix(), // 过期时间为 3 天后
// 		"iat":     time.Now().Unix(),                     // 签发时间
// 		"user_id": id,                                    // 用户ID
// 	}

// 	// 使用声明创建新 Token
// 	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	// 使用配置中的密钥签名 Token
// 	token, err := tk.SignedString([]byte(cf.TokenJWT))
// 	if err != nil {
// 		return "", fmt.Errorf("用户校验失败")
// 	}

// 	return token, nil
// }

// // 验证并解析 Token
// func (_this *Token) ParseToken(tokenStr string) (interface{}, error) {
// 	jws, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("token错误")
// 		}
// 		return []byte(cf.TokenJWT), nil
// 	})

// 	// 检查 Token 是否有效
// 	if err != nil || !jws.Valid {
// 		if errors.Is(err, jwt.ErrTokenMalformed) {
// 			return 0, fmt.Errorf("非法访问")
// 		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
// 			return 0, fmt.Errorf("无效用户")
// 		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
// 			return 0, fmt.Errorf("登录已过期, 请重新登录")
// 		} else {
// 			return 0, err
// 		}
// 	}

// 	claims, ok := jws.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return 0, fmt.Errorf("无效的Token")
// 	}

// 	usrIDFloat, ok := claims["user_id"].(float64)
// 	if !ok {
// 		return 0, fmt.Errorf("无效的Token")
// 	}

// 	return int(usrIDFloat), nil
// }
