package model

type Client struct {
	Rest []RestClient `json:"rest" yaml:"rest"`
}

type ModuleArgs struct {
	Client Client `json:"client" yaml:"client"`
	// LogLevel is the log level, default info
	LogLevel string `json:"log_level" yaml:"log_level"`
}

type ModuleResponse struct {
	Msg    string `json:"msg"`
	Failed bool   `json:"failed"`
}
