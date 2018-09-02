package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gitlab.com/teruhirokomaki/sheetskv/token"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var sheetId string
var sheetName string
var carriageReturn bool

func getKeyList(srv *sheets.Service)  {
	readRange := strings.Replace("'__SHEETNAME__'!A:A", "__SHEETNAME__", sheetName, -1)

	resp, err := srv.Spreadsheets.Values.Get(sheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found")
	} else {
		for _, row := range resp.Values {
			fmt.Println(row[0])
		}
	}
}

func getRow(srv *sheets.Service, key string) (int, int, string) {
	var keyListLen int
	var rowNum int
	var rowValue string
	readRange := strings.Replace("'__SHEETNAME__'!A:B", "__SHEETNAME__", sheetName, -1)
	resp, err := srv.Spreadsheets.Values.Get(sheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	keyListLen = len(resp.Values)
	if keyListLen == 0 {
		fmt.Println("No data found")
	} else {
		for i, row := range resp.Values {
			if key == row[0] {
				if len(row) == 1 { // no value
					rowValue = ""
				} else {
					rowValue, _ = row[1].(string)
				}
				rowNum = i
				break
			}
		}
	}
	return keyListLen, rowNum, rowValue
}

func addRow(srv *sheets.Service, key string, val string) {
	var targetRowNum int

	// getRow
	keyListLen, rowNum, _ := getRow(srv, key)

	// set targetRowNum
	if rowNum == 0 {
		targetRowNum = keyListLen+ 1
	} else {
		targetRowNum = rowNum+ 1
	}

	// updateRange
	updateRangeTemp := "'__SHEETNAME__'!A__ROWNUM__:B__ROWNUM__"
	r := strings.NewReplacer("__SHEETNAME__", sheetName, "__ROWNUM__", strconv.Itoa(targetRowNum))
	updateRange := r.Replace(updateRangeTemp)

	// valueRange
	var row []interface{}
	row = append(row, key)
	row = append(row, val)

	var mapVals [][]interface{}
	mapVals = append(mapVals, row)

	valueRange := &sheets.ValueRange{}
	valueRange.Values = mapVals

	_, err := srv.Spreadsheets.Values.Update(sheetId, updateRange, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
}

func getService() *sheets.Service {
	b, err := ioutil.ReadFile(token.GetCredentialsDir())
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := token.GetClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	return srv
}

func main() {
	app := cli.NewApp()
	app.Name = "sheetskv"
	app.Usage = "sheetskv is a CLI tool for using Google Spreadsheets as the Key Value Store"
	app.Version = "0.1.0"
	app.Author = "frudens Inc. <https://frudens.com>"

	// global option
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "sheetId, i",
			Usage: "Google Spreadsheets id",
			Destination: &sheetId,
		},
		cli.StringFlag{
			Name: "sheetName, n",
			Usage: "Sheet name of Google Spreadsheets",
			Destination: &sheetName,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List contents of column A of Spreadsheets",
			Action:  func(c *cli.Context) error {
				srv := getService()
				getKeyList(srv)
				return nil
			},
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "If the key matches the Spreadsheets' A column, display the contents of column B",
			Action:  func(c *cli.Context) error {
				if c.NArg() == 0  {
					return cli.NewExitError("Key required for argument", 1)
				}
				key := c.Args().Get(0)

				srv := getService()
				_, rowNum, rowValue := getRow(srv, key)
				if rowNum == 0 {
					return cli.NewExitError("Key not found", 2)
				} else if rowValue == "" {
					return cli.NewExitError("Value is not set", 3)
				}
				if carriageReturn {
					fmt.Println(rowValue)
				} else {
					fmt.Print(rowValue)
				}
				return nil
			},

			// get command option
			Flags: []cli.Flag {
				cli.BoolFlag{
					Name: "carriageReturn, cr",
					Usage: "Line break at standard output",
					Destination: &carriageReturn,
				},
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "If the key matches the Spreadsheets' A column, update the contents of column B, and if it does not match, add it",
			Action:  func(c *cli.Context) error {
				if len(c.Args()) < 2 {
					return cli.NewExitError("Key and value are required as arguments", 1)
				}
				key := c.Args().Get(0)
				value := c.Args().Get(1)

				srv := getService()
				addRow(srv, key, value)
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	// app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
