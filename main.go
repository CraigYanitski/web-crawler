package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

type config struct {
    pages               map[string]int
    baseURL             *url.URL
    mu                  *sync.Mutex
    concurrencyControl  chan struct{}
    wg                  *sync.WaitGroup
    maxPages            int
}

func main() {
    args := os.Args[1:]

    if len(args) < 1 {
        fmt.Println("no website provided")
        os.Exit(1)
    } else if len(args) >  3 {
        fmt.Println("too many arguments provided")
        os.Exit(1)
    }

    baseURL, err := url.Parse(args[0])
    if err != nil {
        fmt.Println(err)
        return
    }
    
    // buffer size
    maxConcurrency := 1
    maxPages := 50
    if len(args) > 1 {
        maxConcurrency, err = strconv.Atoi(args[1])
        maxPages, err = strconv.Atoi(args[2])
    }

    fmt.Printf("starting crawl of: %s\n", baseURL.String())
    cfg := config{
        pages:              make(map[string]int),
        baseURL:            baseURL,
        mu:                 &sync.Mutex{},
        concurrencyControl: make(chan struct{}, maxConcurrency),
        wg:                 &sync.WaitGroup{},
        maxPages:           maxPages,
    }

    //cfg.wg.Add(1)
    cfg.crawlPage(args[0])
    cfg.wg.Wait()

    fmt.Printf("Linked Pages (%d)\n", len(cfg.pages))

    if len(cfg.pages) == 0 {
        fmt.Println("no pages detected")
        return
    } else {
        for link, num := range cfg.pages {
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

func (cfg *config) crawlPage(rawCurrentURL string) {
    cfg.wg.Add(1)
    defer cfg.wg.Done()
    cfg.concurrencyControl <- struct{}{}
    defer func() {
        <-cfg.concurrencyControl
    }()
    
    cfg.mu.Lock()
    if len(cfg.pages) >= cfg.maxPages {
        cfg.mu.Unlock()
        return
    }
    cfg.mu.Unlock()

    addPageVisit := func (normalizedURL string) (isFirst bool) {
        isFirst = false
        cfg.mu.Lock()
        defer cfg.mu.Unlock()
        if val, ok := cfg.pages[normalizedURL]; ok {
            cfg.pages[normalizedURL] = val + 1
        } else {
            isFirst = true
            //fmt.Printf("* new domain page %s\n", currentURL)
            cfg.pages[normalizedURL] = 1
        }
        return
    }

    if !strings.Contains(rawCurrentURL, cfg.baseURL.String()) {
        fmt.Printf("%s not in domain of %s\n", rawCurrentURL, cfg.baseURL.String())
        return
    }

    currentURL, err := normalizeURL(rawCurrentURL)
    if err != nil {
        fmt.Println(err)
        return
    }

    //fmt.Printf("updating page count for %s\n", currentURL)
    first := addPageVisit(currentURL)
    if !first {
        return
    }

    fmt.Printf("getting HTML from %s\n", rawCurrentURL)
    currentHTML, err := getHTML(rawCurrentURL)
    if err != nil {
        fmt.Printf("  %s\n", err)
        return
    }

    links, err := getURLsFromHTML(currentHTML, cfg.baseURL.String())
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
        go cfg.crawlPage(links[i])
    }

    return
}
