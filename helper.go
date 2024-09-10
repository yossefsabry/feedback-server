package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

    "github.com/go-playground/validator/v10"
    "github.com/go-playground/universal-translator"
)

// RemoveByIndex removes the item at the specified index from the slice
func RemoveByIndex[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		// Index out of range, return the original slice
		return slice
	}
	// Remove the item at the index
	return append(slice[:index], slice[index+1:]...)
}

// ResponseHelper is a helper function that will convert the slice of products to JSON
func ResponseHelper(w http.ResponseWriter, response Response) error {
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

// CloseBody close the body after the function done
func CloseBody(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		log.Fatalf("erorr found: %v", err)
	}
}

// FormatedErrorValidator is a helper function that will convert the error message to JSON
func FormatedErrorValidator(err error, trans ut.Translator) (errs []error) {
  if err == nil {
    return nil
  }
  validatorErrs := err.(validator.ValidationErrors)
  for _, e := range validatorErrs {
    translatedErr := fmt.Errorf(e.Translate(trans))
    errs = append(errs, translatedErr)
  }
  return errs
}
