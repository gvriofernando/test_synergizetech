package store

import (
	"context"
	"log"
	"strings"

	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
)

type Store struct {
	rd *redis.Client
	db *gorm.DB
}

type Config struct {
	Rd *redis.Client
	Db *gorm.DB
}

func NewStore(cfg Config) Store {
	return Store{
		rd: cfg.Rd,
		db: cfg.Db,
	}
}

func (s Store) GetDetailPlayer(ctx context.Context, playerId string) (res Player, err error) {
	db := s.db
	query := db.Table("player").Where("player_id = ?", playerId).Find(&res)
	if query.Error != nil {
		log.Printf("Error get player detail with id %s : %s", playerId, query.Error)
		return Player{}, query.Error
	}

	return res, nil
}

func (s Store) GetAllPlayers(ctx context.Context, req GetAllPlayersRequest) (res []Player, err error) {
	db := s.db
	query := db.Table("player")

	if req.PlayerId != "" {
		query = query.Where("player_id = ?", req.PlayerId)
	}

	if req.AccountName != "" {
		query = query.Where("account_name = ?", req.AccountName)
	}

	if req.AccountNumber != "" {
		query = query.Where("account_number = ?", req.AccountNumber)
	}

	if req.BankName != "" {
		query = query.Where("bank_name = ?", req.BankName)
	}

	if req.RemainingBalance != 0 {
		query = query.Where("wallet >= ?", req.RemainingBalance)
	}

	query = query.Find(&res)
	if query.Error != nil {
		log.Printf("Error get all player : %s", query.Error)
		return nil, query.Error
	}

	return res, nil
}

func (s Store) AddBankAccount(ctx context.Context, req AddBankAccountRequest) (err error) {
	query := s.db.Table("player")

	query = query.Where("player_id = ?", req.PlayerId).Updates(map[string]interface{}{
		"bank_name":      req.BankName,
		"account_name":   req.AccountName,
		"account_number": req.AccountNumber,
	})

	return nil
}

func (s Store) TopUpWalletRequest(ctx context.Context, playerId string, walletAmount float64) (err error) {
	query := s.db.Table("player")

	query = query.Where("player_id = ?", playerId).Updates(map[string]interface{}{
		"wallet": walletAmount,
	})

	return nil
}

func (s Store) Register(ctx context.Context, req Player) (err error) {
	query := s.db.Table("player")

	query = query.Create(&req)
	if query.Error != nil {
		return query.Error
	}

	return nil
}

func (s Store) Login(ctx context.Context, req LoginRequest) (err error) {
	// Define your hashmap data
	hashmapData := map[string]interface{}{
		"playerId": req.PlayerId,
		"jwtToken": req.JwtToken,
	}

	loginKey := "LOGIN:" + strings.ToUpper(req.PlayerId)
	// Set the hashmap data as the value for the Redis key
	err = s.rd.HMSet(ctx, loginKey, hashmapData).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s Store) CheckExistingDataRedis(ctx context.Context, key string) (res bool, err error) {
	// Set the hashmap data as the value for the Redis key
	existsResults, err := s.rd.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if existsResults == 1 {
		return true, nil
	} else if existsResults == 0 {
		return false, nil
	}

	return false, nil
}

func (s Store) DeleteRedis(ctx context.Context, key string) (err error) {
	// Set the hashmap data as the value for the Redis key
	_, err = s.rd.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}
