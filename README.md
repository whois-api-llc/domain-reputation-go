[![domain-reputation-go license](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![domain-reputation-go made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://pkg.go.dev/github.com/whois-api-llc/domain-reputation-go)
[![domain-reputation-go test](https://github.com/whois-api-llc/domain-reputation-go/workflows/Test/badge.svg)](https://github.com/whois-api-llc/domain-reputation-go/actions/)

# Overview

The client library for
[Domain Reputation API](https://domain-reputation.whoisxmlapi.com/)
in Go language.

The minimum go version is 1.17.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/domain-reputation-go
```

# Examples

Full API documentation available [here](https://domain-reputation.whoisxmlapi.com/api/documentation/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := domainreputation.NewBasicClient(apiKey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := domainreputation.NewClient(apiKey, domainreputation.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

Domain Reputation API lets you calculate a domain's reputation scores based on hundreds of parameters.

```go

// Make request to get the Domain Reputation API response as a model instance.
domainReputationResp, _, err := client.Get(ctx, "whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

log.Println(domainReputationResp.Mode, domainReputationResp.ReputationScore)

// Make request to get raw data in XML.
resp, err := client.GetRaw(context.Background(), "whoisxmlapi.com",
    domainreputation.OptionOutputFormat("XML"))
if err != nil {
    log.Fatal(err)
}

log.Println(string(resp.Body))

```