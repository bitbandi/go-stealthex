package stealthex

import (
	"encoding/json"
	"time"
)

type Trade struct {
	Id             string    `json:"id"`
	Type           string    `json:"type"`
	Timestamp      time.Time `json:"timestamp"`
	CurrencyFrom   string    `json:"currency_from"`
	CurrencyTo     string    `json:"currency_to"`
	AmountFrom     float64   `json:"amount_from,string"`
	ExpectedAmount float64   `json:"expected_amount,string"`
	AmountTo       float64   `json:"amount_to,string"`
	AddressFrom    string    `json:"address_from"`
	AddressTo      string    `json:"address_to"`
	ExtraIdFrom    string    `json:"extra_id_from"`
	ExtraIdTo      string    `json:"extra_id_to"`
	TxFrom         string    `json:"tx_from"`
	TxTo           string    `json:"tx_to"`
	Status         string    `json:"status"`
	RefundAddress  string    `json:"refund_address"`
	RefundExtraId  string    `json:"refund_extra_id"`
}

func (t *Trade) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Trade
	aux := &struct {
		Timestamp string `json:"timestamp"`
		Updated   string `json:"updated_at"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Timestamp, err = time.Parse("2006-01-02T15:04:05.999Z", aux.Timestamp)
	if err != nil {
		return err
	}
	/*
		t.Updated, err = time.Parse("2006-01-02T15:04:05.999Z", aux.Updated)
		if err != nil {
			return err
		}
	*/
	return nil
}
