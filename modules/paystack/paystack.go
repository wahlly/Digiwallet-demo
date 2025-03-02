package paystack

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (p *PaystackClient) InitiateTransaction(amount uint, email string) (map[string]any, error) {
	fmt.Println("burl: ", p.SecretKey)
	url := string(p.BaseUrl) + "/transaction/initialize"
	fmt.Println("url: ", url)
	payload := &InitTxnReqBody{
		Email: email,
		Amount: amount * 100,	//convert to kobo
	}

	//convert the payload to json
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	//create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJson))
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
	var resBody map[string]any
	if err := json.Unmarshal(body, &resBody); err != nil {
		return nil, err
	}

	return resBody, nil
}