package domainreputation

import (
	"net/url"
	"strings"
)

// Option adds parameters to the query.
type Option func(v url.Values)

var _ = []Option{
	OptionOutputFormat("JSON"),
	OptionMode("fast"),
}

// OptionOutputFormat sets Response output format JSON | XML. Default: JSON.
func OptionOutputFormat(outputFormat string) Option {
	return func(v url.Values) {
		v.Set("outputFormat", strings.ToUpper(outputFormat))
	}
}

// OptionMode sets the check mode. API can check your domain in 2 modes: 'fast' - Only select codes will run—i.e.,
// 62 WHOIS Domain status, 82 Malware Databases check, 87 SSL certificate validity, and 93 WHOIS Domain check—while
// other tests and data collectors will be disabled. 'full' - all the data and the tests will be processed.
// Acceptable values: fast | full. Default: fast.
func OptionMode(mode string) Option {
	return func(v url.Values) {
		v.Set("mode", mode)
	}
}
