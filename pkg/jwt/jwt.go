package jwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/loak155/techbranch-backend/pkg/uuid"
)

type JwtManager struct {
	issuer  string
	secret  string
	expires time.Duration
}

type Claims struct {
	jwt.StandardClaims
}

func NewJwtManager(issuer, secret string, expires time.Duration) *JwtManager {
	return &JwtManager{
		issuer:  issuer,
		secret:  secret,
		expires: expires,
	}
}

func (m *JwtManager) GenerateToken(userID int) (token, jti string, err error) {
	jti = uuid.NewUUID()

	claims := Claims{
		jwt.StandardClaims{
			Issuer:    m.issuer,
			Subject:   strconv.Itoa(userID),
			Audience:  m.issuer,
			ExpiresAt: time.Now().Add(m.expires).Unix(),
			NotBefore: time.Now().Add(time.Second * -5).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        jti,
		},
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenObj.SignedString([]byte(m.secret))
	return token, jti, err
}

func (m *JwtManager) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("invalid token")
			}
			return []byte(m.secret), nil
		},
	)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (m *JwtManager) GetExpiresIn() int {
	return int(m.expires.Seconds())
}
