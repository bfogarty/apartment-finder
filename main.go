package main

import (
	"log"
	"time"

	"github.com/docopt/docopt-go"
)

func main() {
	usage := `Apartment Finder

Usage:
  apartment-finder [--no-broker-fee] SITE BUDGET BEDROOMS
  apartment-finder -h | --help

Arguments:
  SITE      The Craigslist subdomain, e.g., "boston"
  BUDGET    The maximum budget, in dollars.
  BEDROOMS  The number of bedrooms.

Options:
  -h --help        Show this screen.
  --no-broker-fee  Find only apartments with no broker's fee.`

	args, _ := docopt.ParseDoc(usage)

	a := &App{}
	err := a.Init("config.yml", args["SITE"].(string))
	if err != nil {
		log.Fatal(err)
	}
	a.Watch(args["BUDGET"].(string), args["BEDROOMS"].(string), args["--no-broker-fee"].(bool), 3*time.Minute)
}
