package main

// Listing represents a Craigslist apartment listing.
type Listing struct {
	Title    string
	URL      string
	Location string
	Price    int64
}

// Btos converts the given boolean to the string "0" or "1", for URL encoding.
func Btos(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
