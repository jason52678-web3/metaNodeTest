package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"time"
)

const TokenExpired = "token_expired"

var mySecret = []byte("MetaNode is coming")

// MyClaims 自定义声明结构体并内嵌jwt.RegisteredClaims
// jwt.RegisteredClaims 包含了标准的JWT声明字段
// 如果想要保存更多信息，可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	// 设置过期时间，从配置中获取小时数并转换为秒
	expireHours := viper.GetInt("auth.jwt_expire")
	expirationTime := time.Now().Add(time.Duration(expireHours) * time.Hour)

	// 创建一个我们自己的声明的数据
	c := MyClaims{
		UserID:   userID,
		Username: username, // 自定义字段
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 过期时间
			Issuer:    "myBlog",                           // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 定义一个接收声明的变量
	var mc = new(MyClaims)

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法是否是我们预期的HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method: " + token.Method.Alg())
		}
		// 返回用于验证的密钥
		return mySecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 校验token是否有效
	if token.Valid {
		return mc, nil
	}

	return nil, errors.New("invalid token")
}
