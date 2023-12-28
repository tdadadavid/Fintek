package api

import (
	"context"
	"database/sql"
	db "github/tdadadavid/fingreat/db/sqlc"
	"github/tdadadavid/fingreat/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Transfer struct {
	server *Server
}

type TransferRequest struct {
	ToAccountID int32 `json:"to_account_id" binding:"required"`
	FromAccountID int32 `json:"from_account_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}


func (t Transfer) router(server *Server) {
	t.server = server;

	transferRouter := server.router.Group("/transfers", AuthenticatedMiddleware());

	transferRouter.POST("", t.makeTransfer);
}


func (t *Transfer) makeTransfer(ctx *gin.Context) {
	userID, err := utils.GetActiveUser(ctx)
	if err != nil {
		return
	}

	transferReq := new(TransferRequest);

	if err := ctx.ShouldBindJSON(&transferReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error() });
		return;
	}

	if transferReq.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot perform transaction." });
		return;
	}

	sourceAccount, err := t.server.queries.GetAccountByID(context.Background(), int64(transferReq.FromAccountID));
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Source account does not exists."})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() })
		return
	}

	if sourceAccount.UserID != int32(userID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Something fishy is going on."});
		return;
	}

	destinationAccount, err := t.server.queries.GetAccountByID(context.Background(), int64(transferReq.ToAccountID));
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Destination account does not exists." })
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() })
		return
	}
	
	if sourceAccount.Currency != destinationAccount.Currency {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Currencies of account do not match"});
		return;
	}



	if sourceAccount.Balance < transferReq.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Inssufficient balance."});
		return;
	}

	transferTxArgs := db.CreateTransferParams{
		FromAccountID: transferReq.FromAccountID,
		ToAccountID: transferReq.ToAccountID,
		Amount: transferReq.Amount,
	}

	transfer, err := t.server.queries.TransferTx(context.Background(), transferTxArgs);
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while making transaction"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"transfer": transfer,
		"message": "Transfer transaction succesfull.",
	})
}