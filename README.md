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

-   **Easy to use**: Get up and running with the library in minutes
-   **Intuitive**: Access your Five9 Domains data using only a few functions
-   **Actively developed**: Ideas and contributions welcomed!

---

<div align="right">

[![Go][golang]][golang-url]
[![Code Coverage][coverage]][coverage-url]
[![Go Reference][goref]][goref-url]
[![Go Report Card][goreport]][goreport-url]

</div>

## Getting Started

For a full API reference, see the official [Five9 REST API documentation](https://webapps.five9.com/assets/files/for_customers/documentation/apis/vcc-agent+supervisor-rest-api-reference-guide.pdf)

### Prerequisites

Download and install Go, version 1.21+, from the [official Go website](https://go.dev/doc/install).

### Install

```shell
go get github.com/equalsgibson/five9-go
```

### Quickstart

To see more detailed examples, checkout the [example](/example/) directory. This will demonstrate how to use the library to make REST requests, or connect to the WebSocket and access the in-memory cache.

#### Lookup all the users within your Five9 Domain

Below is a short example showing how to list all the users within your Five9 Domain using the library.

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

<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to learn, get inspired, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- LICENSE -->

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<!-- CONTACT -->

## Contact

[Chris Gibson (@equalsgibson)](https://github.com/equalsgibson)

Project Link: [https://github.com/equalsgibson/five9-go](https://github.com/equalsgibson/five9-go)

<!-- ACKNOWLEDGMENTS -->

## Acknowledgments

-   Huge thanks to [@aaronellington](https://github.com/aaronellington) for the continued assistance
-   Thanks to Five9 for providing documentation for their API.

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[golang]: https://img.shields.io/badge/v1.21-000?logo=go&logoColor=fff&labelColor=444&color=%2300ADD8
[golang-url]: https://go.dev/
[coverage]: https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fequalsgibson.github.io%2Ffive9-go%2Fcoverage%2Fcoverage.json&query=%24.total&label=Coverage
[coverage-url]: https://equalsgibson.github.io/five9-go/coverage/coverage.html
[goaction]: https://github.com/equalsgibson/five9-go/actions/workflows/go.yml/badge.svg?branch=main
[goaction-url]: https://github.com/equalsgibson/five9-go/actions/workflows/go.yml
[goref]: https://pkg.go.dev/badge/github.com/equalsgibson/five9-go.svg
[goref-url]: https://pkg.go.dev/github.com/equalsgibson/five9-go
[goreport]: https://goreportcard.com/badge/github.com/equalsgibson/five9-go
[goreport-url]: https://goreportcard.com/report/github.com/equalsgibson/five9-go
