package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
    r := strings.NewReader(htmlBody)
    baseURL, err := url.Parse(rawBaseURL)
    fmt.Println(baseURL.String())
    if err != nil {
        return nil, err
    }

    node, err := html.Parse(r)
    if err != nil {
        return nil, err
    }

    urls := checkNodeURLs(node)

    fmt.Println(urls)

    for i, url := range urls {
        if strings.Contains(url, ":") {
            urls[i] = url
        } else {
            urls[i] = baseURL.JoinPath(url).String()
        }
        fmt.Printf("\n%s\n", urls[i])
    }

    return urls, nil
}

func checkNodeURLs(node *html.Node) []string {
    var urls []string
    for child := range node.Descendants() {
        if child.Type == html.ElementNode {
            if child.DataAtom == atom.A {
                for _, a := range child.Attr {
                    if a.Key == "href" {
                        urls = append(urls, a.Val)
                    }
                }
            }
        } else {
            for g := range child.Descendants() {
                urls = checkNodeURLs(g)
            }
        }
    }
    return urls
}
