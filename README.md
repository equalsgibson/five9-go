<!-- markdownlint-configure-file { "MD004": { "style": "consistent" } } -->
<!-- markdownlint-disable MD033 -->

#

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://equalsgibson.github.io/five9-go/logo-dark.png">
    <source media="(prefers-color-scheme: light)" srcset="https://equalsgibson.github.io/five9-go/logo-light.png">
    <img src="https://equalsgibson.github.io/five9-go/logo-light.png" width="410" height="205" alt="Five9Go website">
  </picture>
    <br>
    <strong>Easily integrate your Go application with the Five9 REST and WebSocket API</strong>
</p>

<!-- markdownlint-enable MD033 -->

-----

[![Code Coverage](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fequalsgibson.github.io%2Ffive9-go%2Fcoverage%2Fcoverage.json&query=%24.total&label=Coverage)](https://equalsgibson.github.io/five9-go/coverage/coverage.html)
[![Go](https://github.com/equalsgibson/five9-go/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/equalsgibson/five9-go/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/equalsgibson/five9-go.svg)](https://pkg.go.dev/github.com/equalsgibson/five9-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/equalsgibson/five9-go)](https://goreportcard.com/report/github.com/equalsgibson/five9-go)

## Getting Started

For a full API reference, see the official [Five9 REST API documentation](https://webapps.five9.com/assets/files/for_customers/documentation/apis/vcc-agent+supervisor-rest-api-reference-guide.pdf)

### Install
```shell
go get github.com/equalsgibson/five9-go
```

### Quickstart

To see more detailed examples, checkout the [example](/example/) directory. This will demonstrate how to use the library to connect to the WebSocket and access the in-memory cache.

#### Lookup all the users within your Five9 Domain
```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/aaronellington/zendesk-go/zendesk"
)

func main() {
	ctx := context.Background()

	z := zendesk.NewService(
		os.Getenv("ZENDESK_DEMO_SUBDOMAIN"),
		zendesk.AuthenticationToken{
			Email: os.Getenv("ZENDESK_DEMO_EMAIL"),
			Token: os.Getenv("ZENDESK_DEMO_TOKEN"),
		},
		zendesk.ChatCredentials{
			ClientID:     os.Getenv("ZENDESK_DEMO_CHAT_CLIENT_ID"),
			ClientSecret: os.Getenv("ZENDESK_DEMO_CHAT_CLIENT_SECRET"),
		},
		// Logger is optional, see implementation to see how to add your custom logger here
		zendesk.WithLogger(log.New(os.Stdout, "Zendesk API - ", log.LstdFlags)),
		// Optionally set http.RoundTripper - this is helpful when writing tests
		zendesk.WithRoundTripper(customRoundTripper),
	)

	tags, err := support.Tickets().AddTags(ctx, 6170, zendesk.Tags{
		"foobar",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", tags)
}
```