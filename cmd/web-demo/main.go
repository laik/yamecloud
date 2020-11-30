package main

import (
	"github.com/yametech/yamecloud/pkg/action/api"
)

func main() {
	var app api.Interface
	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
