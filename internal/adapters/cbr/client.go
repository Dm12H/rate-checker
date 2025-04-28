package cbr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sync"

	domain "github.com/Dm12H/rate-checker/internal/domain"
)

type RequestEnvelope struct {
	XMLName xml.Name `xml:"soap12:Envelope"`
	XSI     string   `xml:"xmlns:xsi,attr"`
	XSD     string   `xml:"xmlns:xsd,attr"`
	SOAP    string   `xml:"xmlns:soap12,attr"`
	Body    Body     `xml:"soap12:Body"`
}

type ResponseEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	Payload []byte `xml:",innerxml"`
}

func MarshalRequest(payload any) ([]byte, error) {
	pBytes, err := xml.Marshal(payload)
	if err != nil {
		return nil, err
	}
	envelope := RequestEnvelope{
		XSI:  "http://www.w3.org/2001/XMLSchema-instance",
		XSD:  "http://www.w3.org/2001/XMLSchema",
		SOAP: "http://www.w3.org/2003/05/soap-envelope",
		Body: Body{Payload: pBytes},
	}
	encoded, err := xml.Marshal(envelope)
	return encoded, err
}

func UnmarshalResponse(response []byte, payload any) error {
	var envelope ResponseEnvelope
	err := xml.Unmarshal(response, &envelope)
	if err != nil {
		return err
	}
	innerXML := envelope.Body.Payload
	return xml.Unmarshal(innerXML, payload)
}

func ConsumeEndpoint(cbr *CBRClient, payload any) (responseBytes []byte, err error) {
	encoded, err := MarshalRequest(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling SOAP payload: %v", err)
	}
	byteData := append([]byte(xml.Header), encoded...)
	reader := bytes.NewReader(byteData)
	req, err := http.NewRequest("POST", cbr.Url, reader)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	resp, err := cbr.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()
	responseBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}
	return responseBytes, nil
}

type CBRClient struct {
	Url    string
	Client http.Client
}

func NewCBRClient(url string) CBRClient {
	return CBRClient{Url: url, Client: http.Client{}}
}

func (c *CBRClient) GetCurrencyInfo(request domain.CurDataReq) (domain.ExtCurResp, error) {
	currencyMap, err := EnumValutes(c, request.ToDate)
	if err != nil {
		return domain.ExtCurResp{}, err
	}
	var wg sync.WaitGroup
	response := domain.ExtCurResp{FromDate: request.FromDate, ToDate: request.ToDate}
	for _, code := range request.CurrencyCodes {
		curData := domain.CurInfo{Code: code}
		cbrCode := currencyMap[code]
		response.Data = append(response.Data, &curData)
		if len(cbrCode) == 0 {
			curData.Error = domain.ErrCurCode
			continue
		}
		wg.Add(1)
		go func(cur *domain.CurInfo, code string) {
			defer wg.Done()
			curreq := GetCursDynamicXML{
				FromDate:   request.FromDate.AddDate(0, 0, -15),
				ToDate:     request.ToDate,
				ValutaCode: code,
			}
			results, err := GetCursDynamic(c, curreq)
			if err != nil {
				cur.Error = domain.ErrService
				return
			}
			for _, res := range results {
				point := domain.TimeSeries{
					Date:      res.CursDate,
					Count:     res.Vnom,
					UnitValue: res.VunitRate,
					Value:     res.Vcurs,
				}
				cur.Values = append(cur.Values, point)
			}
		}(&curData, cbrCode)
	}
	wg.Wait()
	return response, nil
}
