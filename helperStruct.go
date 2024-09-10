package main

// Product main struct for the product
type Product struct {
	Id          int    `json:"id" validate:"required,numeric"`
	Name        string `json:"Name" validate:"required,min=3,max=100"`
	Slug        string `json:"Slug" validate:"required,alphanum"`
	Description string `json:"Description" validate:"omitempty,max=500"`
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

// response jwt json response
type responseJwt struct {
	Message string `json:"message"`
}

// Jwks struct 
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}
