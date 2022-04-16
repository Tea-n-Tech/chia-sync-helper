package main

import (
	"fmt"

	"github.com/Tea-n-Tech/chia-sync-helper/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
