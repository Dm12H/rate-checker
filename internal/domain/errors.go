package domain

import "errors"

var ErrCurCode = errors.New("currency code not recognized")
var ErrNoData = errors.New("no data for this currency and time period")
var ErrInvalidDate = errors.New("date range is invalid")
var ErrService = errors.New("error fetching data")
