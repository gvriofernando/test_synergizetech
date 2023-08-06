package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gvriofernando/test_synergizetech/app/core"
	"github.com/gvriofernando/test_synergizetech/app/store"
)

type Handlers struct {
	PCore core.Core
}

// Add or Register a New Player
func (h Handlers) GetAllPlayers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		playerId := ctx.Query("playerId")
		accountName := ctx.Query("accountName")
		accountNumber := ctx.Query("accountNumber")
		bankName := ctx.Query("bankName")
		remainingBalance := ctx.Query("remainingBalance")

		remainingBalanceParseFloat, _ := strconv.ParseFloat(remainingBalance, 64)

		coreRes, err := h.PCore.GetAllPlayers(ctx, store.GetAllPlayersRequest{
			PlayerId:         playerId,
			AccountName:      accountName,
			AccountNumber:    accountNumber,
			BankName:         bankName,
			RemainingBalance: remainingBalanceParseFloat,
		})
		if err != nil {
			ctx.JSON(500, Response{
				ErrorCode: 500,
				Message:   err.Error(),
				Data:      "",
			})
			return
		}

		ctx.JSON(200, Response{
			Message: "success",
			Data:    coreRes,
		})
		return
	}
}

func (h Handlers) GetDetailPlayer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		playerId := ctx.Param("playerId")

		coreRes, err := h.PCore.GetDetailPlayer(ctx, playerId)
		if err != nil {
			ctx.JSON(500, Response{
				ErrorCode: 500,
				Message:   err.Error(),
				Data:      "",
			})
			return
		}

		ctx.JSON(200, Response{
			Message: "success",
			Data:    coreRes,
		})
		return
	}
}

func (h Handlers) AddBankAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req store.AddBankAccountRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, Response{
				ErrorCode: 400,
				Message:   err.Error(),
				Data:      "",
			})
			return
		}

		err := h.PCore.AddBankAccount(ctx, req)
		if err != nil {
			ctx.JSON(500, Response{
				ErrorCode: 500,
				Message:   err.Error(),
				Data:      "",
			})
			return
		}

		ctx.JSON(200, Response{
			Message: "success",
		})
		return
	}
}

func (h Handlers) TopupWallet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req core.TopUpWalletRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, Response{
				ErrorCode: 400,
				Message:   err.Error(),
				Data:      "",
			})
			return
		}

		err := h.PCore.TopUpWallet(ctx, req)
		if err != nil {
			ctx.JSON(500, Response{
				ErrorCode: 500,
				Message:   err.Error(),
				Data:      "",
			})
			return
		}

		ctx.JSON(200, Response{
			Message: "success",
		})
		return
	}
}

func (h Handlers) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req core.RegisterRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, Response{
				ErrorCode: 400,
				Message:   err.Error(),
				Data:      "",
			})
			return
		}

		err := h.PCore.Register(ctx, req)
		if err != nil {
			ctx.JSON(500, Response{
				ErrorCode: 500,
				Message:   err.Error(),
				Data:      "",
			})
			return
		}

		ctx.JSON(200, Response{
			Message: "success",
		})
		return
	}
}

func (h Handlers) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req core.LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, Response{
				ErrorCode: 400,
				Message:   err.Error(),
				Data:      "",
			})
			return
		}

		token, err := h.PCore.Login(ctx, req)
		if err != nil {
			ctx.JSON(500, Response{
				ErrorCode: 500,
				Message:   err.Error(),
				Data:      "",
			})
			return
		}

		ctx.JSON(200, Response{
			Message: "success",
			Data:    token,
		})
		return
	}
}

func (h Handlers) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req core.LogoutRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, Response{
				ErrorCode: 400,
				Message:   err.Error(),
				Data:      "",
			})
			return
		}

		err := h.PCore.Logout(ctx, req)
		if err != nil {
			ctx.JSON(500, Response{
				ErrorCode: 500,
				Message:   err.Error(),
				Data:      "",
			})
			return
		}

		ctx.JSON(200, Response{
			Message: "success",
		})
		return
	}
}
