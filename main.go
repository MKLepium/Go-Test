package main

import (
	"context"
	"fmt"
	"net/http"

	"Yeet/Backend"

	"github.com/eiannone/keyboard"
)

func main() {

	// Initialize the keyboard listener
	if err := keyboard.Open(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer keyboard.Close()
	// Create a context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	// Start web server
	go func(ctx context.Context) {
		fmt.Println("Starting the Server!")
		var srv = Backend.HTTPServer(8080, "localhost")
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Printf("ListenAndServe error: %s\n", err)
			}
		}()
		// Wait for context to be done
		<-ctx.Done()
		srv.Shutdown(ctx)
		fmt.Printf("Server got Killed.\n")
	}(ctx)

	fmt.Printf("Press 'q' to quit.\n")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		// Check if 'q' is pressed
		if char == 'q' || key == keyboard.KeyEsc {
			fmt.Println("Quitting...")
			cancel()
			break
		} else {
			fmt.Println("Invalid input. Press 'q' to quit.")
		}
	}
}
