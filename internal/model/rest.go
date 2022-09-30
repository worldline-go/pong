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

type RestBasicAuth struct {
	// Username is the username
	Username string `json:"username" yaml:"username"`
	// Password is the password
	Password string `json:"password" yaml:"password"`
}

type RestRequest struct {
	// URL could be multiple URLs, separated by space
	URL string `json:"url" yaml:"url"`
	// Method is the HTTP method to use, default is GET
	Method string `json:"method" yaml:"method"`
	// Timeout is in seconds, default 5
	Timeout int `json:"timeout" yaml:"timeout"`
	// Headers is the HTTP headers to be used
	Headers map[string]string `json:"headers" yaml:"headers"`
	// BasicAuth is the basic auth to be used
	BasicAuth *RestBasicAuth `json:"basicAuth" yaml:"basicAuth"`
}

type RestRespond struct {
	// Status is the HTTP status code to be expected
	Status int `json:"status" yaml:"status"`
	// Body is the body to be compared
	Body *RestCheckBody `json:"body" yaml:"body"`
}

type RestCheck struct {
	// Request is the request to be made
	Request RestRequest `json:"request" yaml:"request"`
	// Respond is the response to be expected
	Respond RestRespond `json:"respond" yaml:"respond"`
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
