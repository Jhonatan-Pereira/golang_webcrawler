package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"web_crawler/db"

	"golang.org/x/net/html"
)

type VisitedLink struct {
	Website     string    `bson:"website"`
	Link        string    `bson:"link"`
	VisitedDate time.Time `bson:"visited_date"`
}

var link string

func init() {
	flag.StringVar(&link, "url", "https://aprendagolang.com.br", "url para iniciar visitas")
}

func main() {
	flag.Parse()

	c := make(chan bool)
	go visitLink(link)

	<-c
}

func visitLink(link string) {

	fmt.Printf("Visitando: %s\n", link)

	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[error] status diferente de 200: %d\n", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	extractLinks(doc)
}

func extractLinks(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" {
		// fmt.Println(node.Data)
		for _, attr := range node.Attr {
			// fmt.Println((attr.Key))
			if attr.Key != "href" {
				continue
			}

			link, err := url.Parse(attr.Val)
			if err != nil || link.Scheme == "" || link.Scheme == "mailto" {
				continue
			}

			if db.VisitedLink(link.String()) {
				fmt.Printf("link já visitado: %s\n", link)
				continue
			}

			visitedLink := VisitedLink{
				Website:     link.Host,
				Link:        link.String(),
				VisitedDate: time.Now(),
			}

			db.Insert("links", visitedLink)

			go visitLink(link.String())
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c)
	}
}
