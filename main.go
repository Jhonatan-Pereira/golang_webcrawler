package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"web_crawler/db"
	"web_crawler/website"

	"golang.org/x/net/html"
)

var (
	link   string
	action string
)

func init() {
	flag.StringVar(&link, "url", "https://aprendagolang.com.br", "url para iniciar visitas")
	flag.StringVar(&action, "action", "website", "qual serviço iniciar")
}

func main() {
	flag.Parse()

	switch action {
	case "website":
		website.Run()
	case "webcrawler":
		c := make(chan bool)
		go visitLink(link)
		<-c
	default:
		fmt.Printf("action '%s' não reconhecida", action)
	}

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

			if db.CheckVisitedLink(link.String()) {
				fmt.Printf("link já visitado: %s\n", link)
				continue
			}

			visitedLink := db.VisitedLink{
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
