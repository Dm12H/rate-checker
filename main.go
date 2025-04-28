package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Dm12H/rate-checker/internal/adapters/cbr"
	"github.com/Dm12H/rate-checker/internal/app"
	"github.com/Dm12H/rate-checker/internal/domain"
	"github.com/Dm12H/rate-checker/internal/utils"
)

const CBR_URL = "http://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx"

func processInput(currencies []string, fromDate, toDate time.Time) {
	cbrService := cbr.NewCBRClient(CBR_URL)
	application := app.NewApplication(&cbrService)
	results, err := application.GetCurrencyInfo(domain.CurDataReq{
		CurrencyCodes: currencies,
		FromDate:      fromDate,
		ToDate:        toDate,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
	cursWErrs := make([]int, 0)
	validCurs := make([]int, 0)
	for i, cur := range results.Currencies {
		if cur.Error != nil {
			cursWErrs = append(cursWErrs, i)
		} else {
			validCurs = append(validCurs, i)
		}
	}
	if len(cursWErrs) > 0 {
		fmt.Println("Данные по следующим валютам не были получены:")
		for _, idx := range cursWErrs {
			cur := results.Currencies[idx]
			fmt.Printf("%s: %v\n", cur.Code, cur.Error)
		}
	}
	fmt.Println("-------------------------------")
	if len(validCurs) == 0 {
		return
	}
	maxCur := results.Currencies[results.Max]
	minCur := results.Currencies[results.Min]
	fmt.Printf("Наибольший курс: %d %s = %.4f RUB %s \n", int(maxCur.Count), maxCur.Code, maxCur.Max*maxCur.Count, utils.FormatDate(maxCur.MaxDate))
	fmt.Printf("Наименьший курс: %d %s = %.4f RUB %s \n", int(minCur.Count), minCur.Code, minCur.Min*minCur.Count, utils.FormatDate(minCur.MinDate))
	fmt.Println("-------------------------------")
	fmt.Printf("Средние курсы валют с %s по %s:\n", utils.FormatDate(fromDate), utils.FormatDate(toDate))
	for _, idx := range validCurs {
		cur := results.Currencies[idx]
		fmt.Printf("%d %s = %.4f RUB\n", int(cur.Count), cur.Code, cur.Mean*cur.Count)
	}
}

func main() {
	var currencies []string
	var toDate, fromDate time.Time
	var err error
	currencyStr := flag.String("curr", "", "Comma-separated list of currencies (e.g., USD,EUR,CNY)")
	dateStr := flag.String("date", "", "Date for processing (e.g., 17.04.2024, 2024-04-17, 04/17/2024)")
	helpFlag := flag.Bool("help", false, "Show usage information")
	flag.Parse()
	if *helpFlag {
		flag.Usage() // Print usage information
		os.Exit(0)   // Exit cleanly after showing help
	}
	if *currencyStr == "" {
		currencies = []string{"USD", "EUR", "CNY"}
		fmt.Printf("no currencies set, using default selection: %s\n", strings.Join(currencies, ","))
	} else {
		currencies = utils.ParseCurrencies(*currencyStr)
	}
	if *dateStr == "" {
		toDate = time.Now()
		fmt.Printf("date not set, using current date as the end of 90-day period: %s\n", utils.FormatDate(toDate))
	} else {
		toDate, err = utils.ParseDate(*dateStr)
		if err != nil {
			fmt.Println("Error: Could not parse date. Please select another format, like DD.MM.YYYY")
		}
	}
	fromDate = toDate.AddDate(0, 0, -89)
	processInput(currencies, fromDate, toDate)
}
