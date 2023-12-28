package api

import (
	"context"
	"database/sql"
	"fmt"
	db "github/tdadadavid/fingreat/db/sqlc"
	"github/tdadadavid/fingreat/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Account struct {
	server *Server
}

type AccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

type DepositMoneyRequest struct {
	ToAccountID int64 `json:"to_account_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Reference string `json:"reference" binding:"required"`
}

/////////////////////////////////
// Iniitalizes the router 
// for Auth.
////////////////////////////////
func (a Account) router(server * Server) {
	a.server = server;

	serverGroup := server.router.Group("/accounts", AuthenticatedMiddleware());

	serverGroup.POST("", a.openAccount)
	serverGroup.GET("", a.listUserAccounts)
	serverGroup.POST("/deposit", a.deposit)
}

func (a *Account) openAccount(ctx *gin.Context) {
	userID, err := utils.GetActiveUser(ctx);
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return;
	}

	account := new(AccountRequest);

	if err := ctx.ShouldBindJSON(&account); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}

	args := db.CreateAccountParams{
		UserID: int32(userID),
		Currency: account.Currency,
		Balance: 0,
	}

	newAccount, err := a.server.queries.CreateAccount(context.Background(), args);
	if err != nil {

		if pgError, ok := err.(*pq.Error); ok {
			// postgres unique constraint error code.
			if pgError.Code == "23505" {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Account with this currency already exists.", "status": "Failed"});
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() });
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"account": newAccount,
		"message": fmt.Sprintf("Account with {%v} already created.", account.Currency),
	})

}

func (a *Account) listUserAccounts(ctx *gin.Context) {
	userID, err := utils.GetActiveUser(ctx);
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return;
	}

	accounts, err := a.server.queries.GetAccountByUserID(context.Background(), int32(userID));
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() })
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": accounts,
		"message": "User accounts.",
	})
}

func (a *Account) deposit(ctx *gin.Context) {
	userID, err := utils.GetActiveUser(ctx);
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return;
	}

	depositMoneyReq := new(DepositMoneyRequest);

	if err := ctx.ShouldBindJSON(&depositMoneyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}

	destinationAccount, err := a.server.queries.GetAccountByID(context.Background(), depositMoneyReq.ToAccountID);
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Destination account not found." })
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() })
		return
	}

	if destinationAccount.UserID != int32(userID) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Something fishy is going on." })
		return
	}

	depositMoneyArgs := db.CreateMoneyRecordParams {
		UserID: destinationAccount.UserID,
		Amount: depositMoneyReq.Amount,
		Reference: depositMoneyReq.Reference,
		Status: "pending",
	}

	_, err = a.server.queries.CreateMoneyRecord(context.Background(), depositMoneyArgs);
	if err != nil {
		if pqError, ok := err.(*pq.Error); !ok {
			if pqError.Code == "23505" {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Transaction has been executed previously." });
				return
			}
		}else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() })
			return
		}
	}

	updateUserAccountBalanceArgs := db.UpdateAccountBalanceOptimizedParams {
		Amount: (depositMoneyReq.Amount + destinationAccount.Balance),
		ID: int64(destinationAccount.ID),
	}

	_, err = a.server.queries.UpdateAccountBalanceOptimized(context.Background(), updateUserAccountBalanceArgs);
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error() });
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Deposit successful.",
	})
}