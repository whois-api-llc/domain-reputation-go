package example

import (
	"context"
	"errors"
	domainreputation "github.com/whois-api-llc/domain-reputation-go"
	"log"
)

func GetData(apikey string) {
	client := domainreputation.NewBasicClient(apikey)

	// Get parsed Domain Reputation API response by a domain name as a model instance.
	domainReputationResp, resp, err := client.Get(context.Background(),
		"whoisxmlapi.com",
		// this option is ignored, as the inner parser works with JSON only.
		domainreputation.OptionOutputFormat("XML"))

	if err != nil {
		// Handle error message returned by server.
		var apiErr *domainreputation.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	log.Println(domainReputationResp.Mode, domainReputationResp.ReputationScore)
	// Then print all returned tests and warnings for the domain name.
	for _, obj := range domainReputationResp.TestResults {
		log.Println(obj.Test, obj.TestCode)
		for _, warning := range obj.Warnings {
			log.Println(warning.WarningCode, warning.WarningDescription)
		}
	}

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func GetRawData(apikey string) {
	client := domainreputation.NewBasicClient(apikey)

	// Get raw API response.
	resp, err := client.GetRaw(context.Background(),
		"whoisxmlapi.com",
		// this option causes the all tests to be processed.
		domainreputation.OptionMode("full"),
		domainreputation.OptionOutputFormat("XML"))

	if err != nil {
		// Handle error message returned by server.
		log.Fatal(err)
	}

	log.Println(string(resp.Body))
}
