package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "encoding/json"
)

const uri = "https://api.coindesk.com/v1/bpi/currentprice.json"

type CurrencyResponse struct {
  Code string `json:"code"`
  Symbol string `json:"symbol"`
  Rate string `json:"rate"`
  Description string `json:"description"`
  RateFloat float32 `json:"rate_float"`
}

type BpiResponse struct {
  Usd CurrencyResponse `json:"USD"`
  Eur CurrencyResponse `json:"EUR"`
  Gbp CurrencyResponse `json:"GBP"`
}

type TimeResponse struct {
  Updated string `json:"updated"`
  UpdatedISO string `json:"updatedISO"`
  Updateduk string `json:"updateduk"`
}

type PriceIndexResponse struct {
  Time   TimeResponse `json:"time"`
  Bpi BpiResponse `json:"bpi"`
}

func main() { 
  jsonResponse, jsonErr := http.Get(uri)
  if jsonErr != nil {
  	panic(fmt.Sprintf("can't load price index:", jsonErr))
  }
  defer jsonResponse.Body.Close()

  jsonBody, _ := ioutil.ReadAll(jsonResponse.Body)
  priceIndexResponse := PriceIndexResponse{}
  priceIndexResponseErr := json.Unmarshal(jsonBody, &priceIndexResponse)
  if priceIndexResponseErr != nil {
    fmt.Println("error:", priceIndexResponseErr)
  }
  fmt.Printf("Bitcoin Index: %.3f $\n", priceIndexResponse.Bpi.Usd.RateFloat)
}