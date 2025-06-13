package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
    args := os.Args[1:]

    if len(args) < 1 {
        fmt.Println("no website provided")
        os.Exit(1)
    } else if len(args) > 1 {
        fmt.Println("too many arguments provided")
        os.Exit(1)
    }

    fmt.Printf("starting crawl of: %s\n", args[0])
    
    pages := make(map[string]int)

    crawlPage(args[0], args[0], pages)

    fmt.Printf("Linked Pages (%d)\n", len(pages))

    if len(pages) == 0 {
        fmt.Println("no pages detected")
        return
    } else {
        for link, num := range pages {
            fmt.Printf(" - (%v) %s\n", num, link)
        }
    }

    return
}

func getHTML(rawURL string) (string, error) {
    resp, err := http.Get(rawURL)
    if err != nil {
        return "", err
    }
    if resp.StatusCode > 299 {
        return "", errors.New(fmt.Sprintf("error getting from %s: status code %v", rawURL, resp.StatusCode))
    }
    content := resp.Header.Get("Content-Type")
    if !strings.Contains(content, "text/html") {
        return "", errors.New(fmt.Sprintf("error getting from %s: response content-type %s does not contain text/html", rawURL, content))
    }
    defer resp.Body.Close()

    rawHTML, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(rawHTML), err
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
    if !strings.Contains(rawCurrentURL, rawBaseURL) {
        fmt.Printf("%s not in domain of %s\n", rawCurrentURL, rawBaseURL)
        return
    }

    currentURL, err := normalizeURL(rawCurrentURL)
    if err != nil {
        fmt.Println(err)
        return
    }

    //fmt.Printf("updating page count for %s\n", currentURL)
    if val, ok := pages[currentURL]; ok {
        pages[currentURL] = val + 1
        return
    } else {
        //fmt.Printf("* new domain page %s\n", currentURL)
        pages[currentURL] = 1
    }

    fmt.Printf("getting HTML from %s\n", rawCurrentURL)
    currentHTML, err := getHTML(rawCurrentURL)
    if err != nil {
        fmt.Println(err)
        return
    }

    links, err := getURLsFromHTML(currentHTML, rawBaseURL)
    if err != nil {
        fmt.Println(err)
        return
    }

    for i, link := range links  {
        normLink, err := normalizeURL(link)
        if err != nil {
            fmt.Println(err)
            continue
        }
        if normLink == currentURL {
            //fmt.Printf("skipped: %s same as %s\n", link, rawCurrentURL)
            continue
        }
        crawlPage(rawBaseURL, links[i], pages)
    }

    return
}
