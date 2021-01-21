package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokFile string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	//tokFile := "token.json"
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
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
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

// GetService function
func GetService(credentialFile string, tokFile string) (*drive.Service, error) {
	f, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		log.Fatalf("Unable to read credencial file %s: %v", credentialFile, err)
		return nil, err
	}
	config, err := google.ConfigFromJSON(f, drive.DriveFileScope)
	if err != nil {
		return nil, err
	}
	client := getClient(config, tokFile)
	service, err := drive.New(client)
	if err != nil {
		log.Fatalf("Cannot create the Google Drive service: %v", err)
		return nil, err
	}
	return service, err
}

// CreateDir function
func CreateDir(service *drive.Service, name string, parentID string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentID},
	}
	file, err := service.Files.Create(d).Do()
	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}
	return file, nil
}

// CreateFile function
func CreateFile(service *drive.Service, name string, mimeType string, content io.Reader, parentID string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentID},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}
	return file, nil
}

// DeleteFile function
func DeleteFile(service *drive.Service, fileID string) error {
	err := service.Files.Delete(fileID).Do()
	return err
}

// SearchFile function
func SearchFile(service *drive.Service, q string) (*drive.FileList, error) {
	r, err := service.Files.List().Q(q).Fields("files(id, name)").Do()
	return r, err
}

// ListAllFiles function
func ListAllFiles(service *drive.Service) (*drive.FileList, error) {
	r, err := service.Files.List().Do()
	return r, err
}
