package business

import (
	businessdata "github.com/amir-the-h/okex/models/business"
)

type (
	Candle struct {
		Arg  *map[string]interface{} `json:"arg"`
		Data *businessdata.Candle    `json:"data"`
	}
)
