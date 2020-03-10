package main

import (
	"encoding/json"
	"fmt"
	"github.com/mpurdon/gorest/pkg"
	"os"
	"strings"
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
		Title: "The Title",
		Body: "The post body",
		UserId: 1,
	})

	request := pkg.Request{
		Method:  pkg.Post,
		BaseURL: strings.Join([]string {baseURL, path}, "/"),
		Headers: headers,
		Body:    body,
	}

	response, err := pkg.Send(request)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Status code: %d\n", response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
