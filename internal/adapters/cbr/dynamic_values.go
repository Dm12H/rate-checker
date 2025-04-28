package cbr

import (
	"encoding/xml"
	"fmt"
	"time"
)

type GetCursDynamicXML struct {
	XMLName    struct{}  `xml:"http://web.cbr.ru/ GetCursDynamicXML"`
	FromDate   time.Time `xml:"FromDate"`
	ToDate     time.Time `xml:"ToDate"`
	ValutaCode string    `xml:"ValutaCode"`
}

type GetCursDynamicXMLResponse struct {
	XMLName    xml.Name            `xml:"GetCursDynamicXMLResponse"`
	ValuteData []ValuteCursDynamic `xml:"GetCursDynamicXMLResult>ValuteData>ValuteCursDynamic"`
}

type ValuteCursDynamic struct {
	XMLName   xml.Name  `xml:"ValuteCursDynamic"`
	CursDate  time.Time `xml:"CursDate"`
	Vcode     string    `xml:"Vcode"`
	Vnom      float64   `xml:"Vnom"`
	Vcurs     float64   `xml:"Vcurs"`
	VunitRate float64   `xml:"VunitRate"`
}

func GetCursDynamic(client *CBRClient, request GetCursDynamicXML) ([]ValuteCursDynamic, error) {
	response, err := ConsumeEndpoint(client, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get dynamic data: %v", err)
	}
	var resPayload GetCursDynamicXMLResponse
	err = UnmarshalResponse(response, &resPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return resPayload.ValuteData, nil
}
