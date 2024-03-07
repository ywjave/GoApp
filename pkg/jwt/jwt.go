package jwt

import (
	"errors"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

var mySecret = []byte("jwt加密")

// GenToken ⽣成access token 和 refresh token
// GenToken 生成JWT
func GenToken(userid int64, username string) (string, error) {
	// 创建一个我们自己的声明
	claims := MyClaims{
		userid,
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(),
			Issuer:    "GoApp", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT(accesstoken?)
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	mc := new(MyClaims)
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

//// RefreshToken 刷新AccessToken
//func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
//	// refresh token⽆效直接返回
//	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
//		return
//	}
//	// 从旧access token中解析出claims数据
//	var claims MyClaims
//	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
//	v, _ := err.(*jwt.ValidationError)
//	// 当access token是过期错误 并且 refresh token没有过期时就创建⼀个新的access token
//	if v.Errors == jwt.ValidationErrorExpired {
//		return GenToken(claims.UserID)
//	}
//	return
//}
