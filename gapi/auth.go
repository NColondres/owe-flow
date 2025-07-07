package gapi

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Checks for env variable
func getEnvVar(e string) string {
	v, b := os.LookupEnv(e)
	if b {
		return v
	}
	return ""
}

var Jfile = getEnvVar("SVC_PATH")
var spreadsheetId = getEnvVar("GSHEET")
var people = make(map[string]float32)

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
	// Initialize the Sheets service
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(Jfile), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}
	// Define the spreadsheet ID and range
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
		}
		fmt.Println(people)
	}
}
