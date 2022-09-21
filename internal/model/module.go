package model

type Client struct {
	Rest []RestClient `json:"rest" yaml:"rest"`
}

type ModuleArgs struct {
	Client Client `json:"client" yaml:"client"`
	// LogLevel is the log level, default info
	LogLevel string `json:"log_level" yaml:"log_level"`
	// Delims is the delimeters to use for the template
	Delims []string `json:"delims" yaml:"delims"`
}

type ModuleResponse struct {
	Msg    string `json:"msg"`
	Failed bool   `json:"failed"`
}
