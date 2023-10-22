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
	"net/http"
	"os"

	"github.com/equalsgibson/five9-go/five9"
	"github.com/equalsgibson/five9-go/five9/five9types"
)

func main() {
	// Set up a new Five9 service
	ctx := context.Background()
	c := five9.NewService(
		five9types.PasswordCredentials{
			Username: os.Getenv("FIVE9USERNAME"),
			Password: os.Getenv("FIVE9PASSWORD"),
		},
		five9.AddRequestPreprocessor(func(r *http.Request) error {
			log.Printf("five9 Rest API Call: [%s] %s", r.Method, r.URL.String())

			return nil
		}),
	)

	domainUsers, err := c.Supervisor().GetAllDomainUsers(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Print a count of the users within your Five9 Domain
	log.Printf("You have %d users within your Five9 Domain.\n", len(domainUsers))
}
```