package db

import "context"


type TransferTxResponse struct {
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	EntryIn Entry `json:"entry_in"`
	EntryOut Entry `json:"entry_out"`
	Transfer Transfer `json:"transfer"`
}

func (s *Store) TransferTx(ctx context.Context, TransferReq CreateTransferParams) (TransferTxResponse, error) {
	var tx TransferTxResponse
	var tranErr error

	err := s.executeTransaction(ctx, func(q * Queries) error {
		// transfer money
	
		tx.Transfer, tranErr = q.CreateTransfer(context.Background(), TransferReq)
		if tranErr != nil {
			return tranErr;
		}

		// record entries
			
		// in 
		entryInArgs := CreateEntryParams{
			AccountID: TransferReq.ToAccountID,
			Amount: TransferReq.Amount,
		}
		tx.EntryIn, tranErr = q.CreateEntry(context.Background(), entryInArgs);
		if tranErr != nil {
			return tranErr;
		}
		
		// out
		entryOutArgs := CreateEntryParams{
			AccountID: TransferReq.FromAccountID,
			Amount: -1 * TransferReq.Amount,
		}
		tx.EntryOut, tranErr = q.CreateEntry(context.Background(), entryOutArgs);
		if tranErr != nil {
			return tranErr;
		}

		// update both balance.

		// update recipient balance.
		recipientAccountBalanceArgs := UpdateAccountBalanceOptimizedParams {
			Amount: TransferReq.Amount,
			ID: int64(TransferReq.ToAccountID),
		}
		tx.ToAccount, tranErr = q.UpdateAccountBalanceOptimized(context.Background(), recipientAccountBalanceArgs);
		if tranErr != nil {
			return tranErr;
		}

		// update sender balance.
		senderAccountBalanceArgs := UpdateAccountBalanceOptimizedParams {
			Amount: -1 * TransferReq.Amount,
			ID: int64(TransferReq.FromAccountID),
		}
		tx.FromAccount, tranErr = q.UpdateAccountBalanceOptimized(context.Background(), senderAccountBalanceArgs);
		if tranErr != nil {
			return tranErr;
		}

		return nil;
	});

	return tx, err;
}

