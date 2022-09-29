package model

var DefaultRestCheck = RestCheck{
	Request: RestRequest{
		Method: "GET",
	},
	Respond: RestRespond{
		Status: 200,
	},
}

var DefaultRestClient = RestClient{
	Concurrent: 1,
}
