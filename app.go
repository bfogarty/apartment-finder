package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// App stores the application state.
type App struct {
	// whether or not the app has been initialized
	isInitialized bool
	// the Craigslist site, e.g., "boston"
	site string
	// ID of the newest listing in the current scrape
	currentNewestID string
	// ID of the newest listing in the last scrape
	lastNewestID string
	// listings retrieved in the last scrape
	lastRetrieved int

	config    *Config
	collector *colly.Collector
}

// Init sets up the application state.
func (a *App) Init(configPath string, site string) error {
	if a.isInitialized {
		return ErrAlreadyInitialized
	}

	// set site
	a.site = site

	// load config from file
	config, err := NewConfigFromYAML(configPath)
	if err != nil {
		return err
	}
	a.config = config

	// initialize the collector
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	c.OnRequest(func(_ *colly.Request) {
		a.lastRetrieved = 0
	})
	c.OnHTML("ul.rows", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("li.result-row,h4.nearby", a.handleRow)
	})
	a.collector = c

	return nil
}

// Watch begins monitoring for new listings.
func (a *App) Watch(budget string, bedrooms string) {
	params := url.Values{
		"max_price":        {budget},
		"min_bedrooms":     {bedrooms},
		"max_bedrooms":     {bedrooms},
		"availabilityMode": {"0"},
		"broker_fee":       {"1"},
		"sale_date":        {"all+dates"},
	}.Encode()

	clURL := fmt.Sprintf("https://%s.craigslist.org/search/apa?%s", a.site, params)

	for {
		a.collector.Visit(clURL)
		log.Printf("Retrieved %d listing(s).\n", a.lastRetrieved)
		time.Sleep(3 * time.Second)
	}
}

func (a *App) handleRow(i int, row *colly.HTMLElement) bool {
	id := row.Attr("data-pid")

	if i == 0 {
		a.currentNewestID = id
	}

	// Show only results more recent than the last we've seen, and before the
	// "Few local results found" banner. We rely on the fact that Craigslist sorts
	// the listings newest first, so rather than parse dates ourselves, we break
	// if we see an ID we've seen before.
	if strings.Contains(row.Attr("class"), "nearby") || id == a.lastNewestID {
		a.lastNewestID = a.currentNewestID
		a.currentNewestID = ""
		return false
	}

	title := row.ChildText("p.result-info > a.result-title")
	url := row.ChildAttr("p.result-info > a.result-title", "href")
	loc := row.ChildText("p.result-info > span.result-meta > span.result-hood")
	priceStr := row.ChildText("p.result-info > span.result-meta > span.result-price")

	// trim surrounding parentheses from location
	loc = loc[1 : len(loc)-1]

	// remove leading $ and parse to int
	price, err := strconv.ParseInt(priceStr[1:], 10, 64)
	if err != nil {
		price = -1
	}

	a.notify(&Listing{
		ID:       id,
		Title:    title,
		URL:      url,
		Location: loc,
		Price:    price,
	})

	a.lastRetrieved++
	return true
}

func (a *App) notify(l *Listing) {
	body := fmt.Sprintf("New Listing\nLocation: %s\nRent: %d\n%s", l.Location, l.Price, l.URL)

	msg := url.Values{
		"To":   {a.config.Notifications.RecipientPhone},
		"From": {a.config.Twilio.PhoneFrom},
		"Body": {body},
	}
	msgStr := strings.NewReader(msg.Encode())

	twURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", a.config.Twilio.AccountSID)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", twURL, msgStr)
	req.SetBasicAuth(a.config.Twilio.AccountSID, a.config.Twilio.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client.Do(req)
}
