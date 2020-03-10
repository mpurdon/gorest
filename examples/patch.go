package main

import (
	"encoding/json"
	"fmt"
	"github.com/mpurdon/gorest/pkg"
	"os"
)

// An example of a patch request
func main() {

	headers := make(map[string]string)
	headers["Content-type"] = "application/json; charset=UTF-8"

	apiKey, apiKeyExists := os.LookupEnv("API_KEY")
	if apiKeyExists {
		headers["Authorization"] = "Bearer " + apiKey
	}

	body, err := json.Marshal(Post{
		Title: "The Updated Title",
	})

	request := pkg.Request{
		Method:  pkg.Patch,
		BaseURL: Endpoint(Config.postId),
		Headers: headers,
		Body:    body,
	}

	response, err := pkg.Send(request)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
