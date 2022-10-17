package domainreputation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathDomainReputationResponseOK         = "/DomainReputation/ok"
	pathDomainReputationResponseError      = "/DomainReputation/error"
	pathDomainReputationResponse500        = "/DomainReputation/500"
	pathDomainReputationResponsePartial1   = "/DomainReputation/partial"
	pathDomainReputationResponsePartial2   = "/DomainReputation/partial2"
	pathDomainReputationResponseUnparsable = "/DomainReputation/unparsable"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

// dummyServer is the sample of the Domain Reputation API server for testing.
func dummyServer(resp, respUnparsable string, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathDomainReputationResponseOK:
		case pathDomainReputationResponseError:
			w.WriteHeader(499)
			response = respErr
		case pathDomainReputationResponse500:
			w.WriteHeader(500)
			response = respUnparsable
		case pathDomainReputationResponsePartial1:
			response = response[:len(response)-10]
		case pathDomainReputationResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		case pathDomainReputationResponseUnparsable:
			response = respUnparsable
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

// newAPI returns new Domain Reputation API client for testing.
func newAPI(apiServer *httptest.Server, link string) *Client {
	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}

	apiURL.Path = link

	params := ClientParams{
		HTTPClient:              apiServer.Client(),
		DomainReputationBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

// TestDomainReputationGet tests the Get function.
func TestDomainReputationGet(t *testing.T) {
	checkResultRec := func(res *DomainReputationResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"mode":"fast","reputationScore":98.58,"testResults":[{"test":"SSL vulnerabilities","testCode":88,
"warnings":[{"warningDescription":"HTTP Strict Transport Security not set","warningCode":6015},
{"warningDescription":"TLSA record not configured or configured wrong","warningCode":6019},
{"warningDescription":"OCSP stapling not configured","warningCode":6006}]}]}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory string
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathDomainReputationResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathDomainReputationResponse500,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathDomainReputationResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathDomainReputationResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathDomainReputationResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "API error: [499] Test error message.",
		},
		{
			name: "unparsable response",
			path: pathDomainReputationResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument",
			path: pathDomainReputationResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "domainName" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.Get(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("DomainReputation.Get() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("DomainReputation.Get() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("DomainReputation.Get() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestDomainReputationGetRaw tests the GetRaw function.
func TestDomainReputationGetRaw(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"mode":"fast","reputationScore":98.58,"testResults":[{"test":"SSL vulnerabilities","testCode":88,
"warnings":[{"warningDescription":"HTTP Strict Transport Security not set","warningCode":6015},
{"warningDescription":"TLSA record not configured or configured wrong","warningCode":6019},
{"warningDescription":"OCSP stapling not configured","warningCode":6006}]}]}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory string
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathDomainReputationResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathDomainReputationResponse500,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathDomainReputationResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathDomainReputationResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathDomainReputationResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathDomainReputationResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "API failed with status code: 499",
		},
		{
			name: "invalid argument",
			path: pathDomainReputationResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: `invalid argument: "domainName" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.GetRaw(tt.args.ctx, tt.args.options.mandatory)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("DomainReputation.GetRaw() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("DomainReputation.GetRaw() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}
