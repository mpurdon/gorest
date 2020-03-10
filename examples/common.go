package main

import "strings"

// Config is an anonymous struct to hold configuration details
var Config = struct {
	baseURL	string // The base URL to talk to
	path    string // The URI path
	postId 	string // The Id of the Post
	userId 	string // The Id of the User
}{
	baseURL:	"https://jsonplaceholder.typicode.com",
	path:    	"posts",
	postId:  	"1",
	userId:  	"1",
}

// Endpoint builds the fully-qualified path from a slice of path parts
func Endpoint(pathParts ...string) string {
	return strings.Join(append([]string{Config.baseURL, Config.path}, pathParts...), "/")
}

// Post represents a post by a user
// Generated using https://mholt.github.io/json-to-go/
type Post struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
}


// Post represents a collection of posts by a user
type Posts struct {
	Posts   Post `json:"posts"`
}
