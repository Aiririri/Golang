package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type URLAnalyzer struct{}


func (URLAnalyzer) Validate(input string) error {
	_, err := parseURL(input)
	return err
}

func (URLAnalyzer) Analyze(input string) (Result, error) {
	u, err := parseURL(input)
	if err != nil {
		return Result{}, err
	}

	port := 80
	if u.Scheme == "https" {
		port = 443
	}
	if p := u.Port(); p != "" {
		port, _ = strconv.Atoi(p) 
	}

	path := u.Path
	if path == "" {
		path = "/"
	}

	depth := 0
	for _, segment := range strings.Split(path, "/") {
		if segment != "" {
			depth++
		}
	}

	params := 0
	for _, values := range u.Query() {
		params += len(values)
	}

	return Result{
		Scheme:      u.Scheme,
		Host:        strings.ToLower(u.Hostname()),
		Port:        port,
		IsSecure:    u.Scheme == "https",
		Path:        path,
		PathDepth:   depth,
		QueryParams: params,
		HasFragment: u.Fragment != "",
	}, nil
}

func parseURL(input string) (*url.URL, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return nil, fmt.Errorf("%w: url is empty", ErrInvalidInput)
	}
	u, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("%w: scheme must be http or https, got %q", ErrInvalidInput, u.Scheme)
	}
	if u.Hostname() == "" {
		return nil, fmt.Errorf("%w: host is missing", ErrInvalidInput)
	}
	if p := u.Port(); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 || n > 65535 {
			return nil, fmt.Errorf("%w: port %q must be 1-65535", ErrInvalidInput, p)
		}
	}
	return u, nil
}
