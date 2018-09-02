# sheetskv

sheetskv is a CLI tool for using Google Spreadsheets as the Key Value Store

## Intro

* Use Google Spreadsheets as a key value store.
* Manage by using the column A of the Google Spreadsheets as a key and the column B as the value

## Getting Started

1. You must have a [Go](http://golang.org) compiler installed.
2. Download and build sheetskv: `go get github.com/frudens/sheetskv`
3. Either copy the `sheetskv` executable in `$GOPATH/bin` to a directory in
   your `PATH`, or add `$GOPATH/bin` to your `PATH`.
4. Create a new console project and enable the Google Sheets API.
5. Download the configuration file.
6. Move the downloaded file to your home directory and ensure it is named `.sheetskv.credentials.json`.
7. Check the Spreadsheets id and sheet name in the web browser.
8. Run `sheetskv --sheetId SHEETID --sheetName SHEETNAME ls`
9. For sheet ID and sheet name, set alias to .bashrc. `alias sheetskv='sheetskv -i XXXXX -n default' 
10. Run `sheetskv ls`

### When the command is executed for the first time

* Browse to the provided URL in your web browser.
* Log in with Google account.
* Click the Accept button.
* Copy the code you're given, paste it into the command-line prompt, and press Enter.

## Usage

Download configuration file and move.

```
$ cd
$ mv ~/Downloads/client_secret_111111111111-xxxxxxxxxxxxxxx.apps.googleusercontent.com.jsodn ~/.sheetskv.credentials.json
```

Check the Spreadsheets id and sheet name in the web browser.

**Google Spreadsheets example**

SheetId: `12345`

SheetName: `default`

| A | B |
|:---|:---|
| key1 | value1 |
| key2 | value2 |
| key3 | value3 |

**ls (List contents of column A of Spreadsheets)**

```
~ $ sheetskv -i 12345 -n default ls
key1
key2
key3
```

**get (If the key matches the Spreadsheets' A column, display the contents of column B)**

The default is not to output a line feed.

```
~ $ sheetskv -i 12345 -n default get key1
value1 ~ $
```

Output a line feed with the - cr option.

```
~ $ sheetskv -i 12345 -n default get key1 --cr
value1
~ $
```

**add (If the key matches the Spreadsheets' A column, update the contents of column B, and if it does not match, add it)**

```
~ $ sheetskv -i 12345 -n default add key4 value4
```

Because there is no key, it will be added.

| A | B |
|:---|:---|
| key1 | value1 |
| key2 | value2 |
| key3 | value3 |
| key4 | value4 |

```
~ $ sheetskv -i 12345 -n default add key4 updated
```

Because there is a key, it will be updated.

| A | B |
|:---|:---|
| key1 | value1 |
| key2 | value2 |
| key3 | value3 |
| key4 | updated |

## Author

frudens Inc. <https://frudens.com>

## License

This software is distributed under the
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0),
see LICENSE.txt for more information.
