# apartment-finder

Watches Craigslist for new apartment listings and sends alerts using Twilio SMS.

## Prerequisites

This application requires a Twilio account for sending SMS notifications and a Bitly account for shortening links.

## Installation

### Docker

The easiest way to run the application is using Docker, by pulling the image from GitHub Packages:
```console
$ docker run docker.pkg.github.com/bfogarty/apartment-finder/apartment-finder:latest -h
Apartment Finder

Usage:
  apartment-finder [--no-broker-fee] SITE BUDGET BEDROOMS
  apartment-finder -h | --help

Arguments:
  SITE      The Craigslist subdomain, e.g., "boston"
  BUDGET    The maximum budget, in dollars.
  BEDROOMS  The number of bedrooms.

Options:
  -h --help        Show this screen.
  --no-broker-fee  Find only apartments with no broker's fee.
```

### Building from source

Alternatively, the application can be built from source. Clone the repository, then run:
```console
$ go build .
```

## Configuration

Copy `config.yml.example` to `config.yml` and edit its contents. The configuration can also be specified using environment variables to make deployments easier. Environment variables will supercede values in the config file.

| Key | Description | Required | Environment Variable |
|-------------------------------|-------------------------------------------------------------------------------------------------------------------------------|----------|----------------------|
| twilio.accountSid | Your Twilio account SID. | Yes | `TWILIO_ACCOUNT_SID` |
| twilio.authToken | The auth token for your Twilio account. | Yes | `TWILIO_AUTH_TOKEN` |
| twilio.phoneFrom | The phone number for your Twilio account. Must be in "+10000000000" format. | Yes | `TWILIO_PHONE_FROM` |
| bitly.accessToken | The Generic Access Token for your Bitly account. | Yes | `BITLY_ACCESS_TOKEN` |
| notifications.recipientPhones | A list of recipient phone numbers. Must be in "+10000000000" format. Environment variable may contain a comma-separated list. | Yes | `RECIPIENT_PHONES` |

Your Twilio account SID and auth token can be found in your [Console](https://www.twilio.com/console). Bitly provides [the following documentation](https://dev.bitly.com/v4/#section/Application-using-a-single-account) for obtaining a Generic Access Token.
