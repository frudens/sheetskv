package token

import (
	"os"
	"path"
	"runtime"
)

func GetTokenDir() string {
	return path.Join(HomeDir(), ".sheetskv.token.json")
}

func GetCredentialsDir() string {
	return path.Join(HomeDir(), ".sheetskv.credentials.json")
}
func HomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("APPDATA")
	}
	return os.Getenv("HOME")

}
