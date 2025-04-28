package domain

import (
	"sort"
	"time"

	"github.com/Dm12H/rate-checker/internal/utils"
)

type CurDataReq struct {
	CurrencyCodes []string
	ToDate        time.Time
	FromDate      time.Time
}

type CurDataResp struct {
	Max        int
	Min        int
	Currencies []CurDataEntry
}

type CurDataEntry struct {
	Code  string
	Error error
	Min   float64
	MinDate time.Time
	Max   float64
	MaxDate time.Time
	Mean  float64
	Count float64
}

type ExtCurResp struct {
	ToDate   time.Time
	FromDate time.Time
	Data     []*CurInfo
}

type CurInfo struct {
	Code   string
	Error  error
	Values []TimeSeries
}

func (cur *CurInfo) CheckSorted() bool {
	n := len(cur.Values)
	for i := 0; i < n-1; i++ {
		istep := cur.Values[i].Date
		nextstep := cur.Values[i+1].Date
		if istep.After(nextstep) {
			return false
		}
	}
	return true
}

func (cur *CurInfo) Sort() {
	if cur.CheckSorted() {
		return
	}
	sort.Slice(cur.Values, func(i, j int) bool {
		idate := cur.Values[i].Date
		jdate := cur.Values[j].Date
		return idate.Before(jdate)
	})
}

func (cur *CurInfo) CleanDates() {
	for i := range cur.Values {
		cur.Values[i].Date = utils.CleanDate(cur.Values[i].Date)
	}
}

func (cur *CurInfo) Clean() {
	cur.Sort()
	cur.CleanDates()
}

type TimeSeries struct {
	Date      time.Time
	Count     float64
	UnitValue float64
	Value     float64
}
