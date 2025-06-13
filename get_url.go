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
    if err != nil {
        return nil, err
    }

    node, err := html.Parse(r)
    if err != nil {
        return nil, err
    }

    urls := checkNodeURLs(node)

    for i, url := range urls {
        if strings.Contains(url, ":") {
            continue
        } else {
            urls[i] = baseURL.JoinPath(url).String()
        }
    }

    return urls, nil
}

func checkNodeURLs(node *html.Node) []string {
    var urls []string
    for child := range node.Descendants() {
        fmt.Println(child.Namespace)
        if child.Type == html.ElementNode {
            if child.DataAtom == atom.A {
                for _, a := range child.Attr {
                    if a.Key == "href" {
                        urls = append(urls, a.Val)
                    }
                }
            } else {
                urls = checkNodeURLs(child)
            }
        }
    }
    return urls
}
