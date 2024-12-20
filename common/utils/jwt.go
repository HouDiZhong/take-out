package utils

import (
	"errors"
	"strconv"
	"take-out/config"
	"take-out/global"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomPayload 自定义载荷继承原有接口并附带自己的字段
type CustomPayload struct {
	UserId     uint64
	GrantScope string
	jwt.RegisteredClaims
}

func DeleteRedisToken(uid uint64, secret string) (int64, error) {
	return global.Redis.Del(secret + strconv.FormatUint(uid, 10)).Result()
}

// GenerateToken 生成Token uid 用户id subject 签发对象  secret 加盐
func GenerateToken(uid uint64, jwtOption config.JwtOption) (string, error) {
	ttl, _ := strconv.Atoi(jwtOption.TTL)
	ttlDuration := time.Duration(ttl) * time.Second
	claim := CustomPayload{
		UserId:     uid,
		GrantScope: jwtOption.Secret,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   //签发者
			Subject:   jwtOption.Secret,                                //签发对象
			Audience:  jwt.ClaimStrings{"PC", "Wechat_Program"},        //签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttlDuration)), //过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(jwtOption.Secret))
	global.Redis.Set(jwtOption.Secret+strconv.FormatUint(uid, 10), token, ttlDuration)
	return token, err
}

func ParseToken(token string, secret string) (*CustomPayload, error) {
	// 解析token
	parseToken, err := jwt.ParseWithClaims(token, &CustomPayload{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parseToken.Claims.(*CustomPayload); ok && parseToken.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
