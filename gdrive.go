package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type gDrive struct {
	credentialFile string
	client         *http.Client
}

// newGDrive create new instance of the gDrive structure.
// credentialFile is an input parameter is a name of the file there credential parametere to Google Drive API are stored.
// At the end you've got gDrive with an client for further interaction with Google Drive.
func newGDrive(credentialFile string) (*gDrive, error) {
	// Default value
	if credentialFile == "" {
		credentialFile = "credentials.json"
	}

	gd := &gDrive{
		credentialFile: credentialFile,
	}

	b, err := ioutil.ReadFile(gd.credentialFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}

	if err := gd.getClient(config); err != nil {
		return nil, err
	}

	return gd, nil

}

// Retrieve a token, saves the token, then returns the generated client.
func (gd *gDrive) getClient(config *oauth2.Config) error {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := gd.tokenFromFile(tokFile)
	if err != nil {
		tok, err = gd.getTokenFromWeb(config)
		if err != nil {
			return err
		}
		if err := gd.saveToken(tokFile, tok); err != nil {
			return err
		}
	}
	gd.client = config.Client(context.Background(), tok)
	return nil
}

// Retrieves a token from a local file.
func (gd *gDrive) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func (gd *gDrive) getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve token from web %v", err)
	}
	return tok, nil
}

// Saves a token to a file path.
func (gd *gDrive) saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

// Create new connection to Google Drive for further interaction with Google Drive.
func (gd *gDrive) getService() (*drive.Service, error) {
	return drive.New(gd.client)
}
