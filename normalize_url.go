package main

import (
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        return "", err
    }
    path := strings.TrimRight(parsedURL.Path, "/")
    newURL, err := url.JoinPath(parsedURL.Host, path)
    if err != nil {
        return "", err
    }
    return newURL, nil
}
