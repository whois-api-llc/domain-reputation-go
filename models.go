package domainreputation

import (
	"fmt"
)

// Warning is the warning detected during the test execution.
type Warning struct {
	// WarningDescription is the warning description.
	WarningDescription string `json:"warningDescription"`

	// Warning code is the unique numeric warning codes.
	WarningCode int `json:"warningCode"`
}

// TestResult is a part of the Domain Reputation API response containing the test results.
type TestResult struct {
	// Test is the test name which reduced the final score.
	Test string `json:"test"`

	// TestCode is the unique numeric test identifier.
	TestCode int `json:"testCode"`

	// Warnings is the list of warnings detected during the test execution.
	Warnings []Warning `json:"warnings"`
}

// DomainReputationResponse is a response of Domain Reputation API.
type DomainReputationResponse struct {
	// Mode is the selected test mode.
	Mode string `json:"mode"`

	// ReputationScore is the composite safety score based on numerous security data sources.
	ReputationScore float64 `json:"reputationScore"`

	// TestResults is a part of the Domain Reputation API response containing the test results.
	TestResults []TestResult `json:"testResults"`
}

// ErrorMessage is the error message.
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"messages"`
}

// Error returns error message as a string.
func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("API error: [%d] %s", e.Code, e.Message)
}
