package main

import (
	"net/http"
	"encoding/json"
	"io"
	"log"
)

// removeByIndex removes the item at the specified index from the slice
func RemoveByIndex[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		// Index out of range, return the original slice
		return slice
	}
	// Remove the item at the index
	return append(slice[:index], slice[index+1:]...)
}

// responseHelper is a helper function that will convert the slice of products to JSON
func responseHelper(w http.ResponseWriter, response Response) error {
	// Here we are converting the slice of products to JSON
	payload, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	_, err := w.Write([]byte(payload))
	if err != nil {
		return err
	}
	return nil
}

// close the body after the function done
func closeBody(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
}
