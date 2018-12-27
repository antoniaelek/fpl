package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Printf("%10s %40s %40s\n", "min", "goal", "assist")
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML(".event.popUp.goal > .row:nth-child(2)", func(e *colly.HTMLElement) {
		min := e.ChildText(".min")
		player := e.ChildText(".player")
		if len(player) > 8 && player[len(player)-7:len(player)] == "( pen )" {
			player = player[0 : len(player)-8]
		}
		assist := e.ChildText(".assist")
		if len(assist) > 0 && assist[0:4] == "Ast." {
			assist = assist[5:len(assist)]
		}
		fmt.Printf("%10s %40s %40s\n", min, player, assist)
	})

	c.Visit("https://www.premierleague.com/matchweek/3278/blog")
}
