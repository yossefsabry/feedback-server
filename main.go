package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

    // Endpoint to generate JWT
    r.HandleFunc("/token", AuthPage).Methods("GET")  // for genrate a token
	r.Handle("/products", VerifyJWT(ProductsHandler)).Methods("GET")
	r.Handle("/products/{slug}/info",VerifyJWT(GetFeedbackProduct)).Methods("GET")
	r.Handle("/products/{slug}/add", VerifyJWT(AddFeedbackProduct)).Methods("POST")
	r.Handle("/products/{slug}/update", VerifyJWT(UpdateFeedbackProduct)).Methods("POST")
	r.Handle("/products/delete", VerifyJWT(DeleteFeedbackProduct)).Methods("DELETE")


	// adding the cors wrapper
	corsWrapper := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "Access-Control-Allow-Origin", "Accept", "*"},
		AllowCredentials: true,
	})

	port := os.Getenv("PORT")
	fmt.Printf("serve runing in port %v\n", port)

	if err := http.ListenAndServe(port, corsWrapper.Handler(r)); err != nil {
		log.Fatalf("error happend: %v\n", err)
	}
	return
}

