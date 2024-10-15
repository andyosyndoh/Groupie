package main

import (
	"fmt"
	"groupie/internals/routes"
	"net/http"

	"groupie/utils"
)

func main() {
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

	wrappedMux := routes.RouteChecker(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: wrappedMux,
	}

	fmt.Println("server running @http://localhost:8080\n=====================================")
	err := server.ListenAndServe()
	if err != nil {
		utils.ErrorHandler("web")
	}
}
