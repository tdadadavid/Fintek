package test

import (
	"context"
	db "github/tdadadavid/fingreat/db/sqlc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createRandomAccount(userID int64, t *testing.T) db.Account {
	args := db.CreateAccountParams{
		UserID: int32(userID),
		Currency: "NGN",
		Balance: 200,
	}

	account, err := testQuery.CreateAccount(context.Background(), args);

	assert.NoError(t, err);
	assert.NotEmpty(t, account);

	assert.Equal(t, account.Balance, float64(200));
	assert.Equal(t, account.Currency, "NGN")
	assert.WithinDuration(t, account.CreatedAt, time.Now(), 2 * time.Second)
	
	return account;
}

func TestTransfer(t *testing.T) {	
	sender := createRandomUser();
	reciever := createRandomUser();

	senderAccount := createRandomAccount(sender.ID, t);
	recieverAccount := createRandomAccount(reciever.ID, t);

	transferArgs := db.CreateTransferParams{
		FromAccountID: int32(senderAccount.ID),
		ToAccountID: int32(recieverAccount.ID),
		Amount: 10,
	}
	
	transferRespChannel := make(chan db.TransferTxResponse)
	errorChannel := make(chan error)
	count := 3

	for i := 0; i < count; i++ {
		go func() {
			tx, err := testQuery.TransferTx(context.Background(), transferArgs)

			errorChannel <- err;
			transferRespChannel <- tx;
		}()
	}

	for i := 0; i < count; i++ {
		err := <- errorChannel
		tx := <- transferRespChannel

		assert.NoError(t, err);
		assert.NotEmpty(t, tx)

		// test transfer
		assert.Equal(t, tx.Transfer.FromAccountID, transferArgs.FromAccountID);
		assert.Equal(t, tx.Transfer.ToAccountID, transferArgs.ToAccountID);
		assert.Equal(t, tx.Transfer.Amount, transferArgs.Amount);

		// test entry
		// [EntryIn]
		assert.Equal(t, tx.EntryIn.AccountID, transferArgs.ToAccountID);
		assert.Equal(t, tx.EntryIn.Amount, transferArgs.Amount);

		// [EntryOut]
		assert.Equal(t, tx.EntryOut.AccountID, transferArgs.FromAccountID);
		assert.Equal(t, tx.EntryOut.Amount, -1 * transferArgs.Amount);

	}


	senderAccountFromDB, err := testQuery.GetAccountByID(context.Background(), senderAccount.ID);
	assert.NoError(t, err);
	assert.NotEmpty(t, senderAccountFromDB);


	recieverAccountFromDB, err := testQuery.GetAccountByID(context.Background(), recieverAccount.ID);
	assert.NoError(t, err);
	assert.NotEmpty(t, recieverAccountFromDB);

	newAmount := float64(count * int(transferArgs.Amount))
	assert.Equal(t, senderAccountFromDB.Balance, (senderAccount.Balance - newAmount))
	assert.Equal(t, recieverAccountFromDB.Balance, (senderAccount.Balance + newAmount))

}