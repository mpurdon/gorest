package main

import (
	"fmt"
	"github.com/mpurdon/gorest/pkg"
	"os"
)

// An example of a delete request
func main() {

	headers := make(map[string]string)
	headers["Content-type"] = "application/json; charset=UTF-8"

	apiKey, apiKeyExists := os.LookupEnv("API_KEY")
	if apiKeyExists {
		headers["Authorization"] = "Bearer " + apiKey
	}

	queryParams := make(map[string]string)

	request := pkg.Request{
		Method:      pkg.Delete,
		BaseURL:     Endpoint(Config.postId),
		Headers:     headers,
		QueryParams: queryParams,
	}

	response, err := pkg.Send(request)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Headers)
	}
}
