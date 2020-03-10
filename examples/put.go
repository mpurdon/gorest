package main

import (
	"encoding/json"
	"fmt"
	"github.com/mpurdon/gorest/pkg"
	"os"
)

func main() {
	headers := make(map[string]string)
	headers["Content-type"] = "application/json; charset=UTF-8"

	apiKey, apiKeyExists := os.LookupEnv("API_KEY")
	if apiKeyExists {
		headers["Authorization"] = "Bearer " + apiKey
	}

	body, err := json.Marshal(Post{
		Title: "The Updated Title",
		Body: "The post body",
		UserId: 1,
	})

	request := pkg.Request{
		Method:  pkg.Put,
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
