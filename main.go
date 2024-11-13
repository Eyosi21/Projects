package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Stock struct {
	company, price, change string
}

func main() {
	//The company shorthand names you can add to this array
	tickers := []string{
		"RIVN",
		"SPOT",
		"RKLB",
		"CAVA",
		"SAVE",
	}

	stocks := []Stock{}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("section.yf-k4z9w", func(e *colly.HTMLElement) {
		stock := Stock{}

		// Extract company name
		stock.company = e.ChildText("section.container h1")
		fmt.Println("Company:", stock.company)

		// Extract price
		stock.price = e.ChildText("fin-streamer[data-field='regularMarketPrice'] span")
		fmt.Println("Price:", stock.price)

		// Extract change
		stock.change = e.ChildText("fin-streamer[data-field='regularMarketChangePercent'] span")
		fmt.Println("Change:", stock.change)

		// Append to stocks slice if all fields are not empty
		if stock.company != "" && stock.price != "" && stock.change != "" {
			stocks = append(stocks, stock)
		}
	})

	for _, t := range tickers {
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")
		c.Wait()
	}

	fmt.Println("Stocks collected:", stocks)

	// Write to CSV file
	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	headers := []string{"company", "price", "change"}
	writer.Write(headers)

	for _, stock := range stocks {
		record := []string{stock.company, stock.price, stock.change}
		writer.Write(record)
	}
	writer.Flush()
	fmt.Println("Data saved to stocks.csv")
}
