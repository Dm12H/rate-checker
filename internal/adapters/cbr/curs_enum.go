package cbr

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

type EnumValutesXML struct {
	XMLName struct{} `xml:"http://web.cbr.ru/ EnumValutesXML"`
	Seld    bool     `xml:"Seld"`
}

type EnumValutesXMLResponse struct {
	XMLName    xml.Name     `xml:"http://web.cbr.ru/ EnumValutesXMLResponse"`
	ValuteData []EnumValute `xml:"EnumValutesXMLResult>ValuteData>EnumValutes"`
}

type EnumValute struct {
	XMLName   xml.Name `xml:"EnumValutes"`
	Vcode     string   `xml:"Vcode"`
	VcharCode string   `xml:"VcharCode"`
}

func EnumValutes(client *CBRClient, date time.Time) (map[string]string, error) {
	reqPayload := EnumValutesXML{Seld: false}
	response, err := ConsumeEndpoint(client, reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %v", err)
	}
	var resPayload EnumValutesXMLResponse
	err = UnmarshalResponse(response, &resPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	codes := make(map[string]string)
	for _, valute := range resPayload.ValuteData {
		ISOcode := strings.TrimSpace(valute.VcharCode)
		internalCode := strings.TrimSpace(valute.Vcode)
		codes[ISOcode] = internalCode
	}
	return codes, nil
}
