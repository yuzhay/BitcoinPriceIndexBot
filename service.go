package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type PriceIndexResponse struct {
	Time Time `json:"time"`
	Bpi  Bpi  `json:"bpi"`
}

type Time struct {
	Updated    string `json:"updated"`
	UpdatedISO string `json:"updatedISO"`
	Updateduk  string `json:"updateduk"`
}

type Bpi struct {
	Usd Currency `json:"USD"`
	Eur Currency `json:"EUR"`
	Gbp Currency `json:"GBP"`
}

type Currency struct {
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
	Rate        string  `json:"rate"`
	Description string  `json:"description"`
	RateFloat   float32 `json:"rate_float"`
}

func getBitcoinPriceIndex(url string) (Bpi, error) {
	jsonResponse, err := http.Get(url)
	if err != nil {
		log.Printf("can't load price index:%s", err)
		return Bpi{}, err
	}
	defer jsonResponse.Body.Close()

	jsonBody, _ := ioutil.ReadAll(jsonResponse.Body)
	priceIndexResponse := PriceIndexResponse{}
	err = json.Unmarshal(jsonBody, &priceIndexResponse)
	if err != nil {
		log.Printf("can't unmarshal response:%s", err)
		return Bpi{}, err
	}

	return priceIndexResponse.Bpi, err
}
