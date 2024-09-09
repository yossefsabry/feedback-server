package main

// Product main struct for the product
type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"Name"`
	Slug        string `json:"Slug"`
	Description string `json:"Description"`
}

// Response struct for the response
type Response struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Product interface{} `json:"product"`
}

// FeedbackProduct struct for the feedback product
type FeedbackProduct struct {
	Id          int     `json:"id"`
	 // this pointer gone make the description is optional field
	Description *string `json:"description,omitempty"`
	Name        *string `json:"name,omitempty"`
	Slug        *string `json:"slug,omitempty"`
}

// DeleteProduct struct for the delete product
type DeleteProduct struct {
	Id int `json:"id"`
}

