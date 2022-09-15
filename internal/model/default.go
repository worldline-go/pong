package model

var DefaultRestCheck = RestCheck{
	Method: "GET",
	Status: 200,
}

var DefaultRestClient = RestClient{
	Concurrent: 1,
}
