package main

import (
	"Yeet/frontend"
	"context"
	"fmt"
	"net/http"
)

func main() {
	// Create a context
	ctx := context.Background()
	//ctx, cancel := context.WithCancel(ctx)
	// Start web server
	go func(ctx context.Context) {
		fmt.Println("Hello, World!")
		var srv = frontend.FrontendServer(8080, "localhost")
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Printf("ListenAndServe error: %s\n", err)
			}
		}()
		<-ctx.Done()
		srv.Shutdown(ctx)
		fmt.Printf("Server got Killed.\n")
	}(ctx)

	//cancel()
	select {}

	//

}
