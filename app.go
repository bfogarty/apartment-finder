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
	// ID of the newest listing seen last scrape
	newestIDLastScrape string
	// ID of the newest listing seen this scrape
	newestIDThisScrape string
	// number of new listings seen on the last scrape
	countNewLastScrape int

	bitly     *BitlyClient
	config    *Config
	collector *colly.Collector
}

// Init sets up the application state.
func (a *App) Init(configPath string, site string) error {
	if a.isInitialized {
		return ErrAlreadyInitialized
	}

	// load config from file
	config, err := NewConfigFromYAML(configPath)
	if err != nil {
		return err
	}
	a.config = config

	a.site = site
	a.bitly = NewBitlyClient(a.config.Bitly.AccessToken)

	// initialize the collector
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	c.OnRequest(func(_ *colly.Request) {
		a.reset()
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
		log.Printf(
			"Found %d new listing(s) on the last scrape.\n",
			a.countNewLastScrape,
		)
		time.Sleep(3 * time.Second)
	}
}

func (a *App) isFirstRun() bool {
	return a.newestIDLastScrape == ""
}

func (a *App) shouldBreakAtRow(row *colly.HTMLElement) bool {
	// stop if the row is a "Few local results found" banner
	if strings.Contains(row.Attr("class"), "nearby") {
		return true
	}

	// break immediately on first run or if this is a listing we've seen
	return a.isFirstRun() || row.Attr("data-pid") == a.newestIDLastScrape
}

func (a *App) handleRow(i int, row *colly.HTMLElement) bool {
	if i == 0 {
		a.newestIDThisScrape = row.Attr("data-pid")
	}

	if a.shouldBreakAtRow(row) {
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
		Title:    title,
		URL:      url,
		Location: loc,
		Price:    price,
	})

	a.countNewLastScrape++
	return true
}

func (a *App) reset() {
	a.countNewLastScrape = 0
	a.newestIDLastScrape = a.newestIDThisScrape
	a.newestIDThisScrape = ""
}

func (a *App) notify(l *Listing) {
	shortURL, err := a.bitly.ShortenLink(l.URL)
	if err != nil {
		shortURL = l.URL
	}

	body := fmt.Sprintf(
		"New Listing\nLocation: %s\nRent: %d\n%s",
		l.Location, l.Price, shortURL,
	)

	msg := url.Values{
		"To":   {a.config.Notifications.RecipientPhone},
		"From": {a.config.Twilio.PhoneFrom},
		"Body": {body},
	}
	msgStr := strings.NewReader(msg.Encode())

	endpoint := fmt.Sprintf(
		"https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json",
		a.config.Twilio.AccountSID,
	)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", endpoint, msgStr)
	req.SetBasicAuth(a.config.Twilio.AccountSID, a.config.Twilio.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client.Do(req)
}
