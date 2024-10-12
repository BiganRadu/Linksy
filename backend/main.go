package main

import (
	"fmt"
	"os"
)

func main() {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	fmt.Print(accountName)
}
