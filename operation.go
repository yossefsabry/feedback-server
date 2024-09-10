package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-playground/universal-translator"
	_ "github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10"
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
	err := ResponseHelper(w, response)
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
})

// AddFeedbackProduct adding product feedback
var AddFeedbackProduct = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var (
		body[] byte
		err error
		product Product
		response Response
	)

	// check request body
	if body, err = io.ReadAll(r.Body); err != nil {
		log.Fatalf("error happened: %v", err)
	}

	vars := mux.Vars(r)
	slug := vars["slug"]
	for _, p := range products {
		if p.Slug == slug {
			product = p
		}
	}
	if product.Slug != "" { // check if the product slug Already exists
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("the product already exists"))
		return
	}

	// close the body after the function done
	defer CloseBody(r.Body)

	product.Slug = slug            // adding the slug to the product
	product.Id = len(products) + 1 // adding the id for the product
	if err = json.Unmarshal(body, &product); err != nil { // unmarshal the body to the product struct
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	// starting validator the product
	validator := validator.New()
	err = validator.Struct(product)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate the product: %v", err), http.StatusBadRequest)
		return
	}
	
	// FormatedError := FormatedErrorValidator(err, trans) // using for another way to get the error

	w.Header().Set("content-Type", "application/json")
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
	err = ResponseHelper(w, response)
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
})

// GetFeedbackProduct get the feedback from the feedbacks
var GetFeedbackProduct = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var (
		product Product
		response Response
		err error
	)
	vars := mux.Vars(r)
	slug := vars["slug"]
	for _, p := range products {
		if p.Slug == slug {
			product = p
		}
	}
	w.Header().Set("content-Type", "application/json")

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
	err = ResponseHelper(w, response)
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
})

var UpdateFeedbackProduct = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var (
		body[] byte
		err error
		feedback FeedbackProduct
		updateProduct Product
		productIndex int
	)

	// get the body from the reqeust.Body
	if body, err = io.ReadAll(r.Body); err != nil {
		log.Fatalf("error happened: %v", err)
	}

	// close the body after the function done
	defer CloseBody(r.Body)

	// create the feedback struct to get the values from the body
	vars := mux.Vars(r) // getting the slug from the url
	slug := vars["slug"]
	feedback.Slug = &slug

	// unmarshal the body to the feedback struct
	if err = json.Unmarshal(body, &feedback); err != nil {
		log.Fatalf("error happened: %v", err)
	}

	for i, p := range products { // check for the product found or not
		if p.Id == feedback.Id {
			updateProduct = p
			productIndex = i
		}
	}

	if updateProduct.Slug == "" || updateProduct.Name == "" { // check if the product not found
		response := Response{
			Message: "error happened the product not found",
			Status:  http.StatusNotFound,
			Product: Product{},
		}
		err = ResponseHelper(w, response)
		if err != nil {
			log.Fatalf("erorr found: %v", err)
		}
		return
	} else { // if the product found starting update with the new values
		if feedback.Slug != nil {
			updateProduct.Slug = *feedback.Slug
		}
		if feedback.Name != nil {
			updateProduct.Name = *feedback.Name
		}
		if feedback.Description != nil {
			updateProduct.Description = *feedback.Description
		}
	}

	products = RemoveByIndex(products, productIndex) // remove the old product from the products
	products = append(products, updateProduct) // append the new product to the products
	
	// return the new response for the user
	response := Response{
		Message: "update the product successfully",
		Status: http.StatusOK,
		Product: updateProduct,
	}
	err = ResponseHelper(w, response)
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
	return
})

var DeleteFeedbackProduct = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// get the body from the reqeust.Body
	var (
		body[] byte
		err error
		deleteProduct DeleteProduct
	)
	if body, err = io.ReadAll(r.Body); err != nil {
		log.Fatalf("error happened: %v", err)
	}

	// close the body after the function done
	defer CloseBody(r.Body)

	// create the feedback struct to get the values from the body

	// unmarshal the body to the feedback struct
	if err = json.Unmarshal(body, &deleteProduct); err != nil {
		log.Fatalf("error happened: %v", err)
	}

	if deleteProduct.Id == 0 { // check if the product not found
		// return the new response for the user
		response := Response{
			Message: "can't delete the product with id not found",
			Status: http.StatusNotFound,
			Product: Product{},
		}
		err = ResponseHelper(w, response)
		if err != nil {
			log.Fatalf("erorr found: %v", err)
		}
		return
	}

	for i, p := range products { // check for the product found or not
		if p.Id == deleteProduct.Id {
			products = RemoveByIndex(products, i) // remove the old product from the products
		}
	}
	
	// return the new response for the user
	response := Response{
		Message: "delete the product successfully",
		Status: http.StatusOK,
		Product: Product{},
	}
	err = ResponseHelper(w, response)
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
	return
})
