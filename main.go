package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	urls := os.Args[1:]
	if len(urls) == 0 {
		urls = []string{
			"https://Example.com:8443/docs/go/intro?lang=en&v=1#top",
			"http://example.com",
			"ftp://example.com",
			"",
		}
	}

	m := NewManager(URLAnalyzer{})
	for _, u := range urls {
		t, err := m.Create(u)
		if errors.Is(err, ErrInvalidInput) {
			fmt.Printf("rejected %q: %v\n", u, err)
			continue
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		if err := m.Process(t.ID); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Printf("%s [%s] %s\n  %+v\n", t.ID, t.Status, t.Input, t.Result)
	}
}
