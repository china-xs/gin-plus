// Package jwt_auth
// @author: ekin
// @date: 2022/5/31
// @Description: jwt token解析
package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
)

// var ProviderSet = wire.NewSet(NewJwtAuth, NewOps)
var (
	ErrInvalidToken = errors.New("【jwt_auth】invalid auth token")
	ErrNoClaims     = errors.New("【jwt_auth】no auth params")
)

type (
	JwtOpts struct {
		Secret    string   `yaml:"secret"`
		Path      []string `yaml:"path"`
		Prefix    []string `yaml:"prefix"`
		Whitelist []string `yaml:"whitelist"`
	}

	JwtAuth struct {
		secret        string
		prefixPath    []string
		matchPath     map[string]struct{}
		whitelistPath map[string]struct{}
	}
)

func NewJWTOps(v *viper.Viper) (*JwtOpts, error) {
	var (
		err error
		o   = new(JwtOpts)
	)

	if err = v.UnmarshalKey("jwt", o); err != nil {
		return nil, err
	}
	return o, err
}

func NewJwtAuth(o *JwtOpts) *JwtAuth {
	j := &JwtAuth{
		secret:        "",
		prefixPath:    make([]string, 0),
		matchPath:     make(map[string]struct{}, 0),
		whitelistPath: make(map[string]struct{}, 0),
	}

	return j.SetSecret(o.Secret).
		SetPath(o.Path...).
		SetPrefix(o.Prefix...).
		SetWhitelist(o.Whitelist...)
}

// jwt parser
func newParser() *jwt.Parser {
	return jwt.NewParser(jwt.WithJSONNumber())
}

// setting prefix paths
func (a *JwtAuth) SetPrefix(paths ...string) *JwtAuth {
	a.prefixPath = append(a.prefixPath, paths...)
	return a
}

// setting whitelist paths
func (a *JwtAuth) SetWhitelist(paths ...string) *JwtAuth {
	for _, path := range paths {
		a.whitelistPath[path] = struct{}{}
	}
	return a
}

// setting match paths
func (a *JwtAuth) SetPath(paths ...string) *JwtAuth {
	for _, path := range paths {
		a.matchPath[path] = struct{}{}
	}
	return a
}

// setting secret
func (a *JwtAuth) SetSecret(secret string) *JwtAuth {
	a.secret = secret
	return a
}

// need check sign
func (a *JwtAuth) NeedCheck(path string) bool {

	//whitelist
	if _, exists := a.whitelistPath[path]; exists {
		return false
	}

	//matched path && prefix path
	hasPath := false
	if _, exists := a.matchPath[path]; exists {
		hasPath = true
	} else {
		for _, p := range a.prefixPath {
			if strings.HasPrefix(path, p) {
				hasPath = true
				break
			}
		}
	}
	return hasPath
}

// verifier the token
func (a *JwtAuth) Verifier(r *http.Request) (jwt.MapClaims, error) {
	var token *jwt.Token
	token, err := request.ParseFromRequest(
		r,
		request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(a.secret), nil
		},
		request.WithParser(newParser()))

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrNoClaims
	}

	return claims, nil
}

// create token with map
func (a *JwtAuth) CreateTokenWithMapPayload(payloads map[string]interface{}, seconds int64) (string, error) {
	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = now + seconds //过期时间
	claims["iat"] = now
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(a.secret))
}

func generateJWTSecret(length int) (string, error) {
	// 计算需要生成的随机字节长度
	byteLength := length
	if length%4 != 0 {
		byteLength = length + (4 - length%4)
	}

	// 生成随机字节序列
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 将字节序列转换成 base64 编码的字符串
	secret := base64.URLEncoding.EncodeToString(randomBytes)
	return secret[:length], nil

}
