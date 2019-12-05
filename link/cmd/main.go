package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"os"
)
type Links struct {
	Text string
	Url string
}

func main() {
	file, err := os.Open("./ex3.html")
	if err != nil {
		log.Fatal(err)
	}

	var r io.Reader
	r = file
	doc, err := html.Parse(r)
	links := Links{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links.Url= a.Val
					fmt.Print("URL: ",links.Url)
					break
				}
			}
		}
		if n.Type == html.TextNode && n.Parent.Data=="a"{
			for _, a := range n.Parent.Attr {
				if a.Key == "href" {
					links.Text = n.Data
					fmt.Print(" Text: ", links.Text, "\n")
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}


	f(doc)
}
