package main

type Status string

type GiftCodes struct {
	Updated string  `json:"updated"`
	Codes   []Codes `json:"codes"`
}

type Codes struct {
	Code   string `json:"code"`
	Status Status `json:"status"`
}

type Config struct {
	GiftCodeSourceUrl string `json:"giftcode_source_json"`
	RedeemUrl         string `json:"giftcode_redeem_api_url"`
	PlayerIdsUrl      string `json:"player_ids_csv"`
	EncryptKey        string `json:"sign"`
	KingdomId         string `json:"kingdom"`
	RequestInterval   int    `json:"request_interval"`
}

type RedeemRequest struct {
	Sign      string `json:"sign"`
	PlayerId  string `json:"fid"`
	GiftCode  string `json:"cdk"`
	KingdomId string `json:"kid"`
	Time      int64  `json:"time"`
}

type RedeemResponse struct {
	Message string `json:"msg"`
}
