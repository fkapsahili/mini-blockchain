package types

// TransactionInput represents a reference to a previous transaction output
type TransactionInput struct {
	PrevTxHash  [32]byte // Hash of the previous transaction
	OutputIndex uint32
	PublicKey   []byte
	Signature   []byte
}

// TransactionOutput represents a new output created by a transaction
type TransactionOutput struct {
	Amount        uint64
	PublicKeyHash []byte
}

// Transaction represents a transfer of coins between addresses
type Transaction struct {
	Version  uint32
	Inputs   []TransactionInput
	Outputs  []TransactionOutput
	LockTime uint32 // Earliest time when this transaction can be included
	Hash     [32]byte
}
