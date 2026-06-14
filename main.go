package main

import (
	"context"
	"fmt"
	"os"

	handler "github.com/andrewto30/url-shortener/handler"
)

func main() {
	ctx := context.Background()
	if err := handler.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
