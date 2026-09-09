package main

import (
	"bytes"
	"crypto/md5"
	_ "embed"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

//go:embed config.json
var configString []byte

func main() {

	cfg := getConfig()
	giftCodes, err := getGiftCodes(cfg)

	giftCodes.Codes = GetActive(giftCodes.Codes)

	if err != nil {
		panic(err)
	}

	playerIds, err := getPlayerIds(cfg)

	if err != nil {
		panic(err)
	}

	sendGiftsToPlayers(cfg, giftCodes, playerIds)

}

func getConfig() *Config {
	config := new(Config)
	json.Unmarshal(configString, &config)

	return config
}

func getGiftCodes(c *Config) (giftCodes *GiftCodes, err error) {
	resp, err := http.Get(c.GiftCodeSourceUrl)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &giftCodes)

	if err != nil {
		return nil, err
	}

	return
}

func GetActive(codes []Codes) (active []Codes) {
	for _, code := range codes {
		if code.Status == "active" {
			active = append(active, code)
		}
	}

	return
}

func getPlayerIds(c *Config) (result []string, err error) {
	resp, err := http.Get(c.PlayerIdsUrl)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	r := csv.NewReader(bytes.NewReader(body))
	records, err := r.ReadAll()

	if err != nil {
		return nil, err
	}

	for _, v := range records {
		result = append(result, v...)
	}

	return
}

func sendGiftsToPlayers(config *Config, codes *GiftCodes, playerIds []string) error {
	for _, playerId := range playerIds {
		for _, code := range codes.Codes {

			req := NewRedeemRequest(playerId, code.Code, config.KingdomId)
			req.signWithMd5(config.EncryptKey)

			form := req.Form()

			resp, err := http.Post(config.RedeemUrl, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

			if err != nil {
				return err
			}

			defer resp.Body.Close()
		}
	}

	return nil
}

func (data *RedeemRequest) signWithMd5(encryptKey string) {
	payload := map[string]string{
		"fid":  data.PlayerId,
		"cdk":  data.GiftCode,
		"kid":  data.KingdomId,
		"time": strconv.FormatInt(data.Time, 10),
	}

	sortedKeys := make([]string, 0, len(payload))
	for key := range payload {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	parts := make([]string, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		parts = append(parts, fmt.Sprintf("%s=%s", key, payload[key]))
	}
	encoded := strings.Join(parts, "&")

	sum := md5.Sum([]byte(encoded + encryptKey))

	data.Sign = hex.EncodeToString(sum[:])
}

func (r *RedeemRequest) Form() url.Values {
	return url.Values{
		"sign": {r.Sign},
		"fid":  {r.PlayerId},
		"cdk":  {r.GiftCode},
		"kid":  {r.KingdomId},
		"time": {strconv.FormatInt(r.Time, 10)},
	}
}

func NewRedeemRequest(playerId string, giftCode string, kingdomId string) *RedeemRequest {
	return &RedeemRequest{
		PlayerId:  playerId,
		GiftCode:  "KS0909",
		KingdomId: kingdomId,
		Time:      time.Now().Unix(),
	}
}
