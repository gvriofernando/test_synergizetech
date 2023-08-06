package core

import "github.com/dgrijalva/jwt-go/v4"

type GetPlayerDetailResponse struct {
	PlayerId      string
	BankName      string
	AccountNumber string
	AccountName   string
	Wallet        float64
}

type TopUpWalletRequest struct {
	PlayerId    string  `json:"playerId"`
	TopUpAmount float64 `json:"topUpAmount"`
}

type RegisterRequest struct {
	PlayerId      string  `json:"playerId"`
	BankName      string  `json:"bankName"`
	AccountNumber string  `json:"accountNumber"`
	AccountName   string  `json:"accountName"`
	Wallet        float64 `json:"wallet"`
	Password      string  `json:"password"`
}

type LoginRequest struct {
	PlayerId string `json:"playerId"`
	Password string `json:"password"`
}

type JwtClaims struct {
	jwt.StandardClaims
	Id    string `json:"id"`
	Token string `json:"token"`
}

type LogoutRequest struct {
	PlayerId string `json:"playerId"`
}
