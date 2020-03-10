package main

import (
	"fmt"
	"github.com/mpurdon/gorest/pkg"
	"os"
)

func getMany(request pkg.Request) *pkg.Response {
	fmt.Printf("\n* Getting many posts from: %s\n", request.BaseURL)

	response, err := pkg.Send(request)
	if err != nil {
		fmt.Println(err)
	}

	return response
}

func getOne(request pkg.Request) (Post, error) {
	fmt.Printf("\n* Getting one post entity from: %s\n", request.BaseURL)

	response, err := pkg.Send(request)
	if err != nil {
		fmt.Println(err)
	}

	var p = Post{}
	var _, jsonError = response.Json(&p)

	return p, jsonError
}

func getOneMap(request pkg.Request) (map[string]interface{}, error) {
	fmt.Printf("\n* Getting one post from: %s as a map\n", request.BaseURL)

	response, err := pkg.Send(request)
	if err != nil {
		fmt.Println(err)
	}

	var result, jsonError = response.JsonAsMap()

	return result, jsonError
}


func main()  {

	// Build the request Headers
	headers := make(map[string]string)
	headers["Content-type"] = "application/json; charset=UTF-8"

	apiKey, apiKeyExists := os.LookupEnv("API_KEY")
	if apiKeyExists {
		headers["Authorization"] = "Bearer " + apiKey
	}

	// GET Single Post
	request := pkg.Request{
		Method:  pkg.Get,
		BaseURL: Endpoint(Config.postId),
		Headers: headers,
	}

	p, err := getOne(request)
	if err != nil {
		return
	}

	fmt.Printf("User Id: %d\n", p.UserID)
	fmt.Printf("Title:   %s\n", p.Title)

	pm, err := getOneMap(request)
	if err != nil {
		return
	}

	fmt.Printf("User Id: %0.f\n", pm["userId"])
	fmt.Printf("Title:   %s\n", pm["title"])


	// GET Collection
	request.BaseURL = Endpoint()

	// Build the query parameters
	queryParams := make(map[string]string)
	queryParams["userId"] = Config.userId
	queryParams["limit"] = "3"
	queryParams["offset"] = "0"

	// Add the userId to the query parameters
	request.QueryParams = queryParams

	fmt.Printf("\n\n* Getting posts for user %s from: %s\n\n", Config.userId, request.BaseURL)
	manyResponse := getMany(request)
	fmt.Println(manyResponse.Body)
}
