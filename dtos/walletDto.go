package dtos

type WalletP2PTransferReqDto struct {
	Amount    int64  `json:"amount" validate:"required,min=100"`
	Recipient string `json:"recipient" validate:"required"`
}