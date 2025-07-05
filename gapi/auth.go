package gapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	Jfile = "file"
)

// Function to read the service account json cred file
type SVCaccountKey struct {
	Type       string `json:"type"`
	PID        string `json:"project_id"`
	PKID       string `json:"private_key_id"`
	PK         string `json:"private_key"`
	CEmail     string `json:"client_email"`
	CID        string `json:"client_id"`
	AuthURI    string `json:"auth_uri"`
	TokeURI    string `json:"token_uri"`
	ClientCert string `json:"client_x509_cert_url"`
	UDomain    string `json:"universe_domain"`
}

type Sheet struct {
	SpreadsheetId string `json:"spreadsheetId"`
}

func readJsonFile(fPath string) (SVCaccountKey, error) {
	var user SVCaccountKey
	//file
	file, err := os.Open(fPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return user, err
	}
	defer file.Close()

	byteResult, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return user, err
	}

	err = json.Unmarshal(byteResult, &user)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}
	return user, nil
	//fmt.Printf("User: %+v\n", user)
}

/*
	func main() {
		keyInfo, err := readJsonFile("file")
		if err != nil {
			fmt.Println("Failed to read Json quitting")
			return
		}
		conf := &jwt.Config{
			Email: keyInfo.CEmail,
			// The contents of your RSA private key or your PEM file that contains a private key.
			// If you have a p12 file instead, you can use `openssl` to export the private key into a pem file.
			// The field only supports PEM containers with no passphrase.
			PrivateKey: keyInfo.PK,
			Scopes: []string{
				"https://www.googleapis.com/auth/bigquery",
				"https://www.googleapis.com/auth/blogger",
			},
			TokenURL: google.JWTTokenURL,
			// If you would like to impersonate a user, you can create a transport with a subject.
			// The following GET request will be made on the behalf of user@example.com.
			// Optional.
			Subject: "user@example.com",
		}
		// Initiate an http.Client, the following GET request will be authorized and authenticated on the behalf of user@example.com.
		client := conf.Client(oauth2.NoContext)
		client.Get("...")
	}
*/
func ListSheets() {
	// Path to the service account key file

	// Create a client using the credentials
	ctx := context.Background()
	client, err := sheets.NewService(ctx, option.WithCredentialsFile(Jfile))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Use the client to make API calls
	// Example: List projects
	sheets := client.Spreadsheets
	fmt.Println(sheets)

}

func Balala() {
	ctx := context.Background()

	// Create a new Drive service
	service, err := drive.NewService(ctx, option.WithCredentialsFile(Jfile))
	if err != nil {
		log.Fatalf("Unable to create the drive service: %v", err)
	}

	// List all files with the MIME type of Google Spreadsheet
	query := "mimeType = 'application/vnd.google-apps.spreadsheet'"
	resp, err := service.Files.List().Q(query).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	// Print the names and IDs of the spreadsheets
	for _, file := range resp.Files {
		fmt.Printf("File ID: %s, Name: %s\n", file.Id, file.Name)
	}
}

type Person struct {
	Name  string
	Items map[string]Spent
}

type Spent struct {
	Date   time.Duration
	Amount float64
	Item   string
}

func removeCharacters(input string, characters string) string {
	filter := func(r rune) rune {
		if strings.IndexRune(characters, r) < 0 {
			return r
		}
		return -1
	}
	return strings.Map(filter, input)
}

func ReadSheed() {
	ctx := context.Background()
	people := make(map[string]float32)

	// Initialize the Sheets service
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(Jfile), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	// Define the spreadsheet ID and range
	spreadsheetId := "sheet"
	rangeName := "Sheet1!A1:Z100"

	// Retrieve the data
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, rangeName).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}

	// Process the data
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {

		for _, row := range resp.Values {
			out := fmt.Sprintf("Expecting row to have 4 values. %d given.", len(row))
			fmt.Println(out)
			var name string
			var cost float32
			for i, value := range row {
				switch v := value.(type) {
				case string:
					if i == 2 {
						newstr := removeCharacters(v, "$")
						if s, err := strconv.ParseFloat(newstr, 32); err == nil {
							cost = float32(s)
						}
					}
					if i == 1 {
						name = v
					}
				case float64:
					if i == 2 {
						cost = float32(v)
					}
				case float32:
					if i == 2 {
						cost = v
					}
				case bool:
					fmt.Println("Boolean:", v)
				case nil:
					fmt.Println("Empty cell")
				default:
					fmt.Println("Unknown type:", v)
				}
				if name != "" && cost != 0 {
					if _, ok := people[name]; ok {
						people[name] = people[name] + cost
					} else {
						people[name] = cost
					}
				}

			}
			fmt.Println(people)
		}
	}
}
