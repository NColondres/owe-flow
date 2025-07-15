package gapi

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"reflect"
	"slices"
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

// Harded coded values for now
var (
	Jfile          = getEnvVar("SVC_PATH")
	spreadsheetId  = getEnvVar("GSHEET")
	people         = make(map[string]float32)
	requiredFields = []string{"Date", "Description", "Amount", "Paid By"}
)

type Record struct {
	Date        time.Time
	Description string
	Amount      float32
	PaidBy      string
}

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

// Function that gets the top row of a spreadsheet. Validates that is has the 5 values in it. [Date, Description, Amount, Paid By]
// Throw error if sheet is not to the desired specification.

func validateSheetFields(valueRange *sheets.ValueRange) bool {
	firstRow := valueRange.Values[0]
	slog.Info(fmt.Sprintf("First Row in sheet: %+v\n", firstRow))
	for i, value := range firstRow {
		if i < len(requiredFields) {
			if !slices.Contains(requiredFields, value.(string)) {
				log_string := fmt.Sprintf("%q is not an expected field", value)
				slog.Error(log_string)
				return false
			}
		} else {
			// We don't need to care about any other colums
			// We just want to validate the colums matching requiredFields
			break
		}
	}
	slog.Info("Validation successful")
	return true
}

// This function will take the first row and return a mapping of the requiredFields
// and the position (index) it is in.
func getSheetFieldsMapping(valueRange *sheets.ValueRange) map[string]int {
	fieldMapping := map[string]int{}

	firstRow := valueRange.Values[0]

	for index, value := range firstRow {
		fieldMapping[value.(string)] = index
	}
	return fieldMapping
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
		slog.Error("No data found.")
	} else if !validateSheetFields(resp) {
		slog.Error("Sheet Validation Failed")
	} else {

		topRowFields := getSheetFieldsMapping(resp)
		slog.Info(fmt.Sprintf("Indexes of the sheet's fields: %v\n", topRowFields))
		// Loop through the sheet starting from the second row
		for i, row := range resp.Values[1:] {
			record := Record{}
			for i, value := range row {
				if i < len(requiredFields) {
					switch v := value.(type) {
					case string:
						if i == topRowFields["Amount"] {
							newstr := removeCharacters(v, "$")
							if s, err := strconv.ParseFloat(newstr, 32); err == nil {
								record.Amount = float32(s)
							}
						}
						if i == topRowFields["Paid By"] {
							record.PaidBy = v
						}

						if i == topRowFields["Date"] {
							record.Date, _ = time.Parse("1/2/2006", v)
						}

					case float64:
						if i == topRowFields["Amount"] {
							record.Amount = float32(v)
						}
					case float32:
						if i == topRowFields["Amount"] {
							record.Amount = v
						}
					case bool:
						fmt.Println("Boolean:", v)
					case nil:
						fmt.Println("Empty cell")
					default:
						fmt.Println("Unknown type:", v)
					}
				} else {
					break
				}
			}
			if !reflect.ValueOf(record.PaidBy).IsZero() && !reflect.ValueOf(record.Amount).IsZero() {
				if _, ok := people[record.PaidBy]; ok {
					people[record.PaidBy] = people[record.PaidBy] + record.Amount
				} else {
					people[record.PaidBy] = record.Amount
				}
			}
			slog.Info(fmt.Sprintf("Record %d: %+v\n", i+1, record))
		}
		fmt.Println(people)
		postOwes(srv, calculateCosts())
	}
}

// Function to figure out who paid the most money

func findWhoPaidTheMost() string {
	largestVal := float32(0)
	var doesntOwe string
	for p, s := range people {
		if s > largestVal {
			largestVal = s
			doesntOwe = p
		}
	}
	return doesntOwe
}

// Function to figure out who owes the person that paid the most
func calculateCosts() [][]interface{} {
	bigSpender := findWhoPaidTheMost()
	TotalSpent := float32(0)
	for _, s := range people {
		TotalSpent += s
	}

	e := [][]interface{}{
		{"Person", "Pay Amount", "Pay To"},
	}
	//////////////////////////////////////////
	numberOfPeople := len(people)
	individualShare := TotalSpent / float32(numberOfPeople)
	slog.Info(fmt.Sprintln("Total Spent:", TotalSpent, "Individual Share:", individualShare))
	if numberOfPeople > 1 {
		for name, contribution := range people {
			netAmount := individualShare - contribution
			if netAmount > 0 {
				e = append(e, []interface{}{name, netAmount, bigSpender})
				slog.Info(fmt.Sprintf("%s owes $%.2f to %s\n", name, netAmount, bigSpender))
			} else if netAmount < 0 {
				slog.Info(fmt.Sprintf("%s is owed $%.2f\n", name, -netAmount))
			} else {
				slog.Info(fmt.Sprintf("%s has no balance\n", name))
			}
		}
	} else if numberOfPeople == 1 {
		slog.Info("There is only one person found on the sheet. Nothing to calculate")
	}

	return e
}

// Function to post the amount owed by people
func postOwes(srv *sheets.Service, owed [][]interface{}) {
	writeRange := "Sheet1!H1"
	valueRange := &sheets.ValueRange{
		Values: owed,
	}
	_, err := srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}
}
