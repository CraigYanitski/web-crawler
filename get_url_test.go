package main

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

const (
    testfmt = `
Inputs:
  URL:     %v
  Body:    %v
Expected:  %v
Actual:    %v
`
    testBody = `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`
    )

func TestGetURLsFromHTML(t *testing.T) {
    type testCase struct {
        name      string
        inputURL  string
        inputBody string
        expected  []string
    }

    tests := []testCase{
        {
            name:      "absolute and relative URLs",
            inputURL:  "https://blog.boot.dev",
            inputBody: testBody,
            expected:  []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
        },
    }

    failCount := 0
    passCount := 0

    fmt.Println("\n\nTesting getURLsFromHTML")

    for _, test := range tests {
        fmt.Println("----------------------------------------")
        fmt.Printf("Getting URLs from %v", test.inputURL)

        result, err := getURLsFromHTML(test.inputBody, test.inputURL)
        
        if err != nil {
            t.Errorf("Error: %s", err)
        }

        links := []string{}
        var absLink string

        for _, link := range result {
            if strings.Contains(link, ":") {
                absLink, err = normalizeURL(link)
                if err != nil {
                    t.Errorf("Error: %s", err)
                }
            } else {
                u, err := url.Parse(test.inputURL)
                if err != nil {
                    t.Errorf("Error: %s", err)
                }
                absLink, err = normalizeURL(u.JoinPath(link).String())
                if err != nil {
                    t.Errorf("Error: %s", err)
                }
            }
            links = append(links, absLink)
        }

        if !reflect.DeepEqual(links, test.expected) {
            failCount++
            t.Errorf(testfmt, test.inputURL, test.inputBody, test.expected, result)
        } else {
            passCount++
            fmt.Printf(testfmt, test.inputURL, test.inputBody, test.expected, result)
        }
    }

    fmt.Println("========================================")
    fmt.Printf("%d passed, %d failed\n\n\n", passCount, failCount)
}

