package paystack

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type PaystackClient struct {
	SecretKey	string
	BaseUrl	string
}

func NewPaystackClient() *PaystackClient {
	return &PaystackClient{
		SecretKey: os.Getenv("PAYSTACK_SK"),
		BaseUrl: os.Getenv("PAYSTACK_BASEURL"),
	}
}

type InitTxnReqBody struct {
	Email string `json:"email"`
	Amount uint `json:"amount"`
}

type PaystackInitTxnRes struct {
	Status bool `json:"status"`
	Message string `json:"message"`
	Data map[string]any `json:"data"`
}

func (p *PaystackClient) InitiateTransaction(amount uint, email string, ref string) (*PaystackInitTxnRes, error) {
	url := string(p.BaseUrl) + "/transaction/initialize"
	var resBody *PaystackInitTxnRes

	// payload := &InitTxnReqBody{
	// 	Email: email,
	// 	Amount: amount * 100,	//convert to kobo
	// }

	payload := map[string]any{
		"email": email,
		"amount": amount * 100,
		"reference": ref,
	}

	//convert the payload to json
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return resBody, err
	}

	//create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJson))
	if err != nil {
		return resBody, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.SecretKey)

	//send request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return resBody, err
	}
	defer res.Body.Close()

	//read the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return resBody, err
	}

	//parse the response
	if err := json.Unmarshal(body, &resBody); err != nil {
		return resBody, err
	}

	return resBody, nil
}

type PaystackVerifyTxnRes struct {
	Status bool `json:"status"`
	Message string `json:"message"`
	Data map[string]any `json:"data"`
}

func (p *PaystackClient) VerifyTransaction(uid uint, reference string) (*PaystackVerifyTxnRes, error) {
	url := string(p.BaseUrl) + "/transaction/verify/"+reference
	resBody := &PaystackVerifyTxnRes{}

	//create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.SecretKey)

	//send request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	//read the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//parse the response
	if err := json.Unmarshal(body, &resBody); err != nil {
		return nil, err
	}

	return resBody, nil
}

func (p *PaystackClient) CalculateProcessingFee(amount uint) uint{
	var processingCharges uint
	if amount < 2500 {
		processingCharges = uint(0.015 * float64(amount))
	} else{
		processingCharges = uint(0.015 * float64(amount)) + 100
	}

	return processingCharges
}