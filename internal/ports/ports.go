package ports

import "github.com/Dm12H/rate-checker/internal/domain"

type ApplicationPort interface {
	GetCurrencyInfo(domain.CurDataReq) (domain.CurDataResp, error)
}

type ExtCurResourcePort interface {
	GetCurrencyInfo(domain.CurDataReq) (domain.ExtCurResp, error)
}
