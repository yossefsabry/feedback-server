package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt"
)
var sampleSecretKey = []byte("welcome from yossef")

func GenerateJWT() (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)
    
    claims := token.Claims.(jwt.MapClaims)
    claims["exp"] = time.Now().Add(10 * time.Minute).Unix() // expires after 10 minutes
    claims["authorized"] = true
    claims["user"] = "username"

    tokenString, err := token.SignedString(sampleSecretKey)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

// AuthPage generates a JWT token and sends it to the client
func AuthPage(writer http.ResponseWriter, request *http.Request) {
    token, err := GenerateJWT()
    if err != nil {
        http.Error(writer, "Error generating token: "+err.Error(), http.StatusInternalServerError)
        return
    }
    response := map[string]string{"token": token}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(response)
}



// Sample secret key used for signing (HS256)

// Mock function to retrieve the public key for token validation
func getPublicKey(token *jwt.Token) (interface{}, error) {
    // Return the secret key for HS256
    return sampleSecretKey, nil
}

// VerifyJWT middleware function
func VerifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
    return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        // Retrieve the token from the Authorization header
        tokenString := request.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(writer, "Token missing", http.StatusUnauthorized)
            return
        }

        // Remove "Bearer " prefix if present
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Check if the token's signing method is what you expect
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            // Return the secret key used to verify the token
            return getPublicKey(token)
        })

        if err != nil {
            http.Error(writer, "Invalid token: "+err.Error(), http.StatusUnauthorized)
            return
        }

        // Check if the token is valid
        if !token.Valid {
            http.Error(writer, "Token is not valid", http.StatusUnauthorized)
            return
        }

        // Call the original endpoint handler if the token is valid
        endpointHandler(writer, request)
    })
}

