package main

import (
	"fmt"
	"reflect"
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
    testBodyOne = `
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
    testBodyTwo = `
<html>
    <body>
        <a href="https://other.com/path/one">
            <span>Boot.dev</span>
        </a>
        <ul>
            <li>
                <a href="/path/one">
                    <span>Boot.dev</span>
                </a>
            </li>
            <li>
                <a href="/path/two">
                    <span>Boot.dev</span>
                </a>
            </li>
        </ul>
    </body>
</html>
`
    testBodyThree = `
<html>
    <body>
        <a href="http://other.com/path/one">
            <span>Boot.dev</span>
        </a>
        <ul>
            <li>
                <a href="/path/one">
                    <span>Boot.dev</span>
                </a>
            </li>
            <li>
                <a href="path/two">
                    <span>Boot.dev</span>
                </a>
            </li>
        </ul>
        <a href="mailto:boot@dev">
            <span>boot@dev</span>
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
            inputBody: testBodyOne,
            expected:  []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
        }, {
            name:      "nested anchors",
            inputURL:  "https://blog.boot.dev",
            inputBody: testBodyTwo,
            expected:  []string{"https://other.com/path/one", "https://blog.boot.dev/path/one", "https://blog.boot.dev/path/two"},
        }, {
            name:      "alternative schemes",
            inputURL:  "https://blog.boot.dev",
            inputBody: testBodyThree,
            expected:  []string{"http://other.com/path/one", "https://blog.boot.dev/path/one", "https://blog.boot.dev/path/two", "mailto:boot@dev"},
        },
    }

    failCount := 0
    passCount := 0

    fmt.Println("\n\nTesting getURLsFromHTML")

    for _, test := range tests {
        fmt.Println("----------------------------------------")
        fmt.Printf("%v\n", test.name)
        fmt.Printf("Getting URLs from %v\n", test.inputURL)

        result, err := getURLsFromHTML(test.inputBody, test.inputURL)
        
        if err != nil {
            t.Errorf("Error: %s", err)
        }

        // links := []string{}
        // var absLink string

        // for _, link := range result {
        //     if strings.Contains(link, ":") {
        //         absLink, err = normalizeURL(link)
        //         if err != nil {
        //             t.Errorf("Error: %s", err)
        //         }
        //     } else {
        //         u, err := url.Parse(test.inputURL)
        //         if err != nil {
        //             t.Errorf("Error: %s", err)
        //         }
        //         absLink, err = normalizeURL(u.JoinPath(link).String())
        //         if err != nil {
        //             t.Errorf("Error: %s", err)
        //         }
        //     }
        //     links = append(links, absLink)
        // }

        if !reflect.DeepEqual(result, test.expected) {
            failCount++
            t.Errorf(testfmt, test.inputURL, test.inputBody, test.expected, result)
        } else {
            passCount++
            //fmt.Printf(testfmt, test.inputURL, test.inputBody, test.expected, result)
        }
    }

    fmt.Println("========================================")
    fmt.Printf("%d passed, %d failed\n\n\n", passCount, failCount)
}

