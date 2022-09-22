package model

type RestCheckBodyVarsFrom struct {
	// Query get the value from the query string
	Query []string `json:"query" yaml:"query"`
}

type RestCheckBodyVars struct {
	// Set is the set of variables
	Set map[string]interface{} `json:"set" yaml:"set"`
	// From is the source of the variable
	From RestCheckBodyVarsFrom `json:"from" yaml:"from"`
}

type RestCheckBody struct {
	// Variable hold the variables to be used in the template
	Variable RestCheckBodyVars `json:"variable" yaml:"variable"`
	// Raw is the raw body to be compared
	Raw *string `json:"raw" yaml:"raw"`
	// Map is the body to be compared
	Map *string `json:"map" yaml:"map"`
}

type RestCheck struct {
	// URL could be multiple URLs, separated by space
	URL string `json:"url" yaml:"url"`
	// Method is the HTTP method to use, default is GET
	Method string `json:"method" yaml:"method"`
	// Status to check, default 200
	Status int `json:"status" yaml:"status"`
	// Body to check
	Body *RestCheckBody `json:"body" yaml:"body"`
	// Timeout is in seconds, default 5
	Timeout int `json:"timeout" yaml:"timeout"`
}

type RestSetting struct {
	InsecureSkipVerify bool `json:"insecureSkipVerify" yaml:"insecureSkipVerify"`
}

type RestClient struct {
	// Concurrent is the number of concurrent requests, default 1
	Concurrent int `json:"concurrent" yaml:"concurrent"`
	// Setting is the setting for the client
	Setting RestSetting `json:"setting" yaml:"setting"`
	// Checks is the list of checks
	Check []RestCheck `json:"check" yaml:"check"`
}
