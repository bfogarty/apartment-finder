package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Listing represents a Craigslist apartment listing.
type Listing struct {
	Title    string
	Location string
	Price    int64
}

func main() {
	args := os.Args
	if len(os.Args) != 4 {
		fmt.Println("Usage:", os.Args[0], "<site>", "<budget>", "<bedrooms>")
		return
	}
	site := args[1]
	budget := args[2]
	bedrooms := args[3]

	c := colly.NewCollector()
	listings := make([]Listing, 0, 50)

	c.OnHTML("ul.rows", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("li.result-row,h4.nearby", func(_ int, row *colly.HTMLElement) bool {
			// show only results before the "Few local results found" banner
			if strings.Contains(row.Attr("class"), "nearby") {
				return false
			}

			title := row.ChildText("p.result-info > a.result-title")
			loc := row.ChildText("p.result-info > span.result-meta > span.result-hood")
			priceStr := row.ChildText("p.result-info > span.result-meta > span.result-price")

			// trim surrounding parentheses from location
			loc = loc[1 : len(loc)-1]

			// remove leading $ and parse to int
			price, err := strconv.ParseInt(priceStr[1:], 10, 64)
			if err != nil {
				price = -1
			}

			l := Listing{
				Title:    title,
				Location: loc,
				Price:    price,
			}

			listings = append(listings, l)
			return true
		})
	})

	params := url.Values{
		"max_price":        {budget},
		"min_bedrooms":     {bedrooms},
		"max_bedrooms":     {bedrooms},
		"availabilityMode": {"0"},
		"broker_fee":       {"1"},
		"sale_date":        {"all+dates"},
	}.Encode()

	c.Visit(fmt.Sprintf("https://%s.craigslist.org/search/apa?%s", site, params))

	// print the output
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(listings)
}
