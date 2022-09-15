package model

type RestCheck struct {
	// URL could be multiple URLs, separated by space
	URL string `json:"url" yaml:"url"`
	// Method is the HTTP method to use, default is GET
	Method string `json:"method" yaml:"method"`
	// Status to check, default 200
	Status int `json:"status" yaml:"status"`
	// Timeout is in seconds, default 5
	Timeout int `json:"timeout" yaml:"timeout"`
}

type RestClient struct {
	// Concurrent is the number of concurrent requests, default 1
	Concurrent int         `json:"concurrent" yaml:"concurrent"`
	Check      []RestCheck `json:"check" yaml:"check"`
}

type Check struct {
	Rest []RestClient `json:"rest" yaml:"rest"`
}

type ModuleArgs struct {
	Check Check `json:"check" yaml:"check"`
	// LogLevel is the log level, default info
	LogLevel string `json:"log_level" yaml:"log_level"`
}

type ModuleResponse struct {
	Msg    string `json:"msg"`
	Failed bool   `json:"failed"`
}
