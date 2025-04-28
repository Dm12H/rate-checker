package app

import (
	"math"
	"time"

	"github.com/Dm12H/rate-checker/internal/domain"
	"github.com/Dm12H/rate-checker/internal/ports"
)

type Application struct {
	ExternalService ports.ExtCurResourcePort
}

func NewApplication(service ports.ExtCurResourcePort) Application {
	return Application{ExternalService: service}
}

func (app Application) GetCurrencyInfo(request domain.CurDataReq) (domain.CurDataResp, error) {
	response := domain.CurDataResp{}
	startDate, endDate := request.FromDate, request.ToDate
	if endDate.Before(startDate) {
		return response, domain.ErrInvalidDate
	}
	extData, err := app.ExternalService.GetCurrencyInfo(request)
	if err != nil {
		return response, domain.ErrService
	}
	for _, cur := range extData.Data {
		cur.Clean()
		entry := calcMean(cur, request.FromDate, request.ToDate)
		response.Currencies = append(response.Currencies, entry)
	}
	var idxMax, idxMin int
	var maxVal float64 = 0
	minVal := math.MaxFloat64
	for i, cur := range response.Currencies {
		if cur.Error == nil {
			if cur.Max > maxVal {
				maxVal = cur.Max
				idxMax = i
			}
			if cur.Min < minVal {
				minVal = cur.Min
				idxMin = i
			}
		}
	}
	response.Max = idxMax
	response.Min = idxMin
	return response, nil
}

func calcMean(cur *domain.CurInfo, start, end time.Time) domain.CurDataEntry {
	result := domain.CurDataEntry{Code: cur.Code}
	if cur.Error != nil {
		result.Error = cur.Error
		return result
	}
	if len(cur.Values) == 0 {
		result.Error = domain.ErrNoData
		return result
	}
	totalDays := tdiff(start, end) + 1
	var accum, lastval, maxval float64
	var minDate, maxDate time.Time
	minval := math.MaxFloat64
	lastdate := start
	for _, entry := range cur.Values {
		date := entry.Date
		lastval = entry.UnitValue
		if date.Before(start) {
			continue
		}
		if lastval < minval{
			minval = lastval
			minDate = lastdate
		}
		if lastval > maxval{
			maxval = lastval
			maxDate = lastdate
		}
		days := tdiff(lastdate, date)
		accum += lastval * days
		lastdate = date
	}
	rem := tdiff(lastdate, end) + 1
	accum += min(totalDays, rem) * lastval
	result.Mean = accum / totalDays
	result.Min = minval
	result.Max = maxval
	result.MaxDate = maxDate
	result.MinDate = minDate
	result.Count = cur.Values[0].Count
	return result
}

func tdiff(startTime, endTime time.Time) float64 {
	return endTime.Sub(startTime).Hours() / 24
}
