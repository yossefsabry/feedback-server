package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

// Product main struct for the product
type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"Name"`
	Slug        string `json:"Slug"`
	Description string `json:"Description"`
}

type Response struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Product interface{} `json:"product"`
}

type FeedbackProduct struct {
	Id          int     `json:"id"`
	Description *string `json:"description,omitempty"` // this pointer gone make the description is optional field
	Name        *string `json:"name,omitempty"`
	Slug        *string `json:"slug,omitempty"`
}

var (
	products = []Product{
		{Id: 1, Name: "Latte", Slug: "latte", Description: "Frothy milky coffee"},
		{Id: 2, Name: "Espresso", Slug: "espresso", Description: "Short and strong coffee without milk"},
		{Id: 3, Name: "Americano", Slug: "americano", Description: "Black coffee"},
		{Id: 4, Name: "Cappuccino", Slug: "cappuccino", Description: "Coffee with frothy milk"},
		{Id: 5, Name: "Mocha", Slug: "mocha", Description: "Cappuccino with chocolate"},
		{Id: 6, Name: "Macchiato", Slug: "macchiato", Description: "Espresso with a little milk"},
		{Id: 7, Name: "Flat White", Slug: "flat-white", Description: "Coffee with milk and microfoam"},
		{Id: 8, Name: "Long Black", Slug: "long-black", Description: "Double espresso with water"},
	}
)

// NotImplemented the function for handle the response from the user
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Not implemented welcome welcome from yossef"))
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
})

// StatusHandler the function for the status end point for the server
var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("API is up and running"))
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
})

// ProductsHandler getting all products for user
var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Here we are converting the slice of products to JSON
	response := Response{
		Message: "get all products successfully",
		Status:  http.StatusOK,
		Product: products,
	}
	payload, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	_, err := w.Write([]byte(payload))
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
})

// AddFeedbackProduct adding product feedback
var AddFeedbackProduct = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var product Product
	vars := mux.Vars(r)
	slug := vars["slug"]
	// fmt.Printf("vars: %v, and slug: %v\n", vars, slug)
	for _, p := range products {
		if p.Slug == slug {
			product = p
		}
	}
	// Read and print the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("erorr found: %v", err)
		}
	}(r.Body)

	product.Slug = slug            // adding the slug to the product
	product.Id = len(products) + 1 // adding the id for the product
	if err := json.Unmarshal(body, &product); err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}
	w.Header().Set("content-Type", "application/json")
	var response Response
	if product.Slug != "" {
		response = Response{
			Message: "successful adding the product",
			Status:  http.StatusOK,
			Product: product,
		}
		products = append(products, product)
	} else {
		response = Response{
			Message: "error happened : product not found",
			Status:  http.StatusNotFound,
			Product: Product{},
		}
	}
	// convert the struct of response to json
	payload, _ := json.Marshal(response)
	w.WriteHeader(response.Status)
	_, err = w.Write([]byte(payload))
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
})

// GetFeedbackProduct get the feedback from the feedbacks
var GetFeedbackProduct = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var product Product
	vars := mux.Vars(r)
	slug := vars["slug"]
	for _, p := range products {
		if p.Slug == slug {
			product = p
		}
	}
	w.Header().Set("content-Type", "application/json")

	var response Response
	if product.Slug != "" {
		response = Response{
			Message: "successfully get the product",
			Status:  http.StatusOK,
			Product: product,
		}
	} else {
		response = Response{
			Message: "error happened the slug must hava a value",
			Status:  http.StatusNotFound,
			Product: Product{},
		}
	}
	// convert the struct of response to json
	payload, _ := json.Marshal(response)
	w.WriteHeader(response.Status)
	_, err := w.Write([]byte(payload))
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
})

var UpdateFeedbackProduct = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("error happened: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("erorr found: %v", err)
		}
	}(r.Body)
	var feedback FeedbackProduct
	vars := mux.Vars(r) // getting the slug from the url
	slug := vars["slug"]
	feedback.Slug = &slug
	// fmt.Printf("body: %v", string(body))
	if err := json.Unmarshal(body, &feedback); err != nil {
		log.Fatalf("error happened: %v", err)
	}

	var updateProduct *Product
	for _, p := range products {
		if p.Id == feedback.Id {
			updateProduct = &p
		}
	}
	if feedback.Slug != nil {
		updateProduct.Slug = *feedback.Slug
	}
	if feedback.Name != nil {
		updateProduct.Name = *feedback.Name
	}
	if feedback.Description != nil {
		updateProduct.Description = *feedback.Description
	}

	fmt.Printf("the prdouct: %v", feedback)
	w.WriteHeader(200)
	_, err = w.Write([]byte("welcome from yossef"))
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
})
