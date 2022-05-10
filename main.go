package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
	"web_crawler/db"

	"golang.org/x/net/html"
)

var (
	visited map[string]bool = map[string]bool{}
)

type VisitedLink struct {
	Website     string    `bson:"website"`
	Link        string    `bson:"link"`
	VisitedDate time.Time `bson:"visited_date"`
}

func main() {
	visitLink("https://aprendagolang.com.br")

	// fmt.Println(len(links))
}

func visitLink(link string) {
	if ok := visited[link]; ok {
		return
	}

	visited[link] = true

	fmt.Println(link)

	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("status diferente de 200: %d", resp.StatusCode))
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
			if err != nil || link.Scheme == "" {
				continue
			}

			// fmt.Println(link.String())
			// links = append(links, link.String())

			visitedLink := VisitedLink{
				Website:     link.Host,
				Link:        link.String(),
				VisitedDate: time.Now(),
			}

			db.Insert("links", visitedLink)

			visitLink(link.String())
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c)
	}
}
