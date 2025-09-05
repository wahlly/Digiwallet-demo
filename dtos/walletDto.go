package dtos

type WalletP2PTransferReqDto struct {
	Amount    int64  `json:"amount" validate:"required,min=100"`
	Recipient string `json:"recipient" validate:"required"`
}

type InitTxnReqBody struct {
	Email  string `json:"email" validate:"required"`
	Amount int64  `json:"amount" validate:"required"`
}

type PaystackInitTxnRes struct {
	Status  bool           `json:"status" validate:"required"`
	Message string         `json:"message" validate:"required"`
	Data    map[string]any `json:"data" validate:"required"`
}
