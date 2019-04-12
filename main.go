package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// getGroupsForUser returns all groups that start with the string "AWS" for a user
func getGroupsForUser(client admin.Service, user admin.User) ([]string, error) {
	groups, err := client.Groups.List().UserKey(user.Id).Do()
	if err != nil {
		return nil, err
	}
	var result []string
	for _, g := range groups.Groups {
		if strings.HasPrefix(g.Name, "AWS") {
			result = append(result, g.Name)
		}
	}
	return result, nil
}

// getCustomSchemaForUser prepares the custom schema attributes for a user
func getCustomSchemaForUser(client admin.Service, user admin.User) ([]map[string]string, error) {
	groups, err := getGroupsForUser(client, user)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve groups for user: %v", err)
	}
	var result []map[string]string
	for _, g := range groups {
		x := strings.Split(g, "-")[1:]
		accountId := x[0]
		role := x[1]
		result = append(result, map[string]string{
			"type":  "work",
			"value": fmt.Sprintf("arn:aws:iam::%[0]v:role/%[1]v,arn:aws:iam::%[0]v:saml-provider/Google", accountId, role)})
	}
	return result, nil

}

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, admin.AdminDirectoryUserReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := admin.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve directory Client %v", err)
	}

	r, err := srv.Users.List().Domain("superluminar.io").OrderBy("email").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve users in domain: %v", err)
	}

	for _, u := range r.Users {
	}
}
