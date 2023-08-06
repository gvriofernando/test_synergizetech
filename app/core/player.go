package core

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
	"github.com/gvriofernando/test_synergizetech/app/store"

	"github.com/dgrijalva/jwt-go/v4"
)

type Core struct {
	store store.Store
}

type Config struct {
	Rd *redis.Client
	Db *gorm.DB
}

func NewCore(cfg Config) Core {
	return Core{
		store: store.NewStore(store.Config{
			Rd: cfg.Rd,
			Db: cfg.Db,
		}),
	}
}

func (c Core) GetDetailPlayer(ctx context.Context, playerId string) (res GetPlayerDetailResponse, err error) {
	storeRes, err := c.store.GetDetailPlayer(ctx, playerId)

	res = GetPlayerDetailResponse{
		PlayerId:      storeRes.PlayerId,
		BankName:      storeRes.BankName,
		AccountNumber: storeRes.AccountNumber,
		AccountName:   storeRes.AccountName,
		Wallet:        storeRes.Wallet,
	}
	if err != nil {
		return GetPlayerDetailResponse{}, nil
	}

	return res, nil
}

func (c Core) GetAllPlayers(ctx context.Context, req store.GetAllPlayersRequest) (res []GetPlayerDetailResponse, err error) {
	storeRes, err := c.store.GetAllPlayers(ctx, req)

	for key := range storeRes {
		res = append(res, GetPlayerDetailResponse{
			PlayerId:      storeRes[key].PlayerId,
			BankName:      storeRes[key].BankName,
			AccountNumber: storeRes[key].AccountNumber,
			AccountName:   storeRes[key].AccountName,
			Wallet:        storeRes[key].Wallet,
		})
	}

	if err != nil {
		return []GetPlayerDetailResponse{}, nil
	}

	return res, nil
}

func (c Core) AddBankAccount(ctx context.Context, req store.AddBankAccountRequest) (err error) {
	err = c.store.AddBankAccount(ctx, req)

	if err != nil {
		return err
	}

	return nil
}

func (c Core) TopUpWallet(ctx context.Context, req TopUpWalletRequest) (err error) {
	playerData, err := c.GetDetailPlayer(ctx, req.PlayerId)
	if err != nil {
		return err
	}

	currentAmount := playerData.Wallet + req.TopUpAmount

	err = c.store.TopUpWalletRequest(ctx, req.PlayerId, currentAmount)
	if err != nil {
		return err
	}

	return nil
}

func (c Core) Register(ctx context.Context, req RegisterRequest) (err error) {
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	err = c.store.Register(ctx, store.Player{
		PlayerId:       req.PlayerId,
		BankName:       req.BankName,
		AccountNumber:  req.AccountNumber,
		AccountName:    req.AccountName,
		Wallet:         req.Wallet,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c Core) Login(ctx context.Context, req LoginRequest) (token string, err error) {
	storeRes, err := c.store.GetDetailPlayer(ctx, req.PlayerId)
	if err != nil {
		return "", err
	}

	if result := checkPasswordHash(ctx, req.Password, storeRes.HashedPassword); !result {
		err = errors.New("Password didn't match")
		return "", err
	}

	signedToken, err := createJwtToken(ctx, req.PlayerId)
	if err != nil {
		err = errors.New("Error creating JWT Token")
		return "", err
	}

	err = c.store.Login(ctx, store.LoginRequest{
		PlayerId: req.PlayerId,
		JwtToken: signedToken,
	})
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (c Core) Logout(ctx context.Context, req LogoutRequest) (err error) {
	keys := "LOGIN:" + strings.ToUpper(req.PlayerId)

	//check keys
	checkExist, err := c.store.CheckExistingDataRedis(ctx, keys)
	if err != nil {
		return err
	}

	if checkExist {
		//delete from redis
		err = c.store.DeleteRedis(ctx, keys)
		if err != nil {
			return err
		}
	} else {
		err = errors.New("Player hasn't Login")
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

func checkPasswordHash(ctx context.Context, password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createJwtToken(ctx context.Context, playerId string) (jwtToken string, err error) {
	secretKey := []byte(os.Getenv("JWT_SIGNATURE_KEY"))
	loginExpInMinute, _ := strconv.Atoi(os.Getenv("JWT_LOGIN_EXP_IN_MINUTE"))
	loginExp := time.Duration(loginExpInMinute) * time.Minute
	signedMethod := jwt.SigningMethodHS256

	//* jwt payload
	claims := JwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(loginExp)),
		},
		Id:    playerId,
		Token: hash(playerId + time.Now().String()),
	}

	token := jwt.NewWithClaims(signedMethod, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return signedToken, nil
}

func hash(data string) string {
	stringData := []byte(data + time.Now().String())
	hashed := md5.Sum(stringData)
	return hex.EncodeToString(hashed[:])
}
