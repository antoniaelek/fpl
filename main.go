package main

import (
	"fmt"
	"fpl/scraper"
)

func main() {
	scores := scraper.Scrape(3278)
	fmt.Println(len(scores))
}
