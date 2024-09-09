package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	r := mux.NewRouter()
	//  on the default page we will serve  our static index page....
	r.Handle("/", http.FileServer(http.Dir("./views")))
	// we will set up our server so we can serve static assets like images , css files
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("/static"))))

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error happend: %v\n", err)
	}

	////////////
	// NEW CODE
	////////////
	// our api is going to consist of there routes
	// /status - which we will al to make sure that our api up and running
	// /products - which will receiver  a list of products  that the user can leave feedback on
	// /products/{slug}/feedback - which will capture user feedback on products
	// starting creating the end points for the api
	r.Handle("/status", StatusHandler).Methods("GET")
	r.Handle("/products", ProductsHandler).Methods("GET")
	r.Handle("/products/{slug}/info", GetFeedbackProduct).Methods("GET")
	r.Handle("/products/{slug}/add", AddFeedbackProduct).Methods("POST")
	r.Handle("/products/{slug}/update", UpdateFeedbackProduct).Methods("POST")

	port := os.Getenv("PORT")
	fmt.Printf("serve runing in port %v\n", port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("error happend: %v\n", err)
	}
	return
}
