package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"github.com/form3tech-oss/jwt-go"
	"github.com/auth0/go-jwt-middleware"

)


func getPemCert(token *jwt.Token) (string, error) {
    cert := ""
    resp, err := http.Get("http://localhost:8080/.well-known/jwks.json")
    if err != nil { return cert, err }
    defer resp.Body.Close()

    var jwks Jwks
    err = json.NewDecoder(resp.Body).Decode(&jwks)
    if err != nil { return cert, err }

    kid, ok := token.Header["kid"].(string)
    if !ok {
        return cert, errors.New("missing or invalid 'kid' in token header")
    }

    for _, key := range jwks.Keys {
        if kid == key.Kid {
            cert = "-----BEGIN CERTIFICATE-----\n" + key.X5c[0] + "\n-----END CERTIFICATE-----"
            break
        }
    }

    if cert == "" {
        return cert, errors.New("unable to find appropriate key in JWKS")
    }
    return cert, nil
}

// Define the JWT Middleware
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
        // Verify audience claim
        aud := "http://localhost:8080"  // Updated for local testing
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok { return nil, errors.New("invalid token claims") }
        if !claims.VerifyAudience(aud, false) {
            return nil, errors.New("invalid audience")
        }

        // Verify issuer claim
        iss := "http://localhost:8080"  // Updated for local testing
        if !claims.VerifyIssuer(iss, false) {
            return nil, errors.New("invalid issuer")
        }

        // Get the PEM certificate
        cert, err := getPemCert(token)
        if err != nil { return nil, err }

        // Parse the public key from PEM
        publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
        if err != nil { return nil, err }
        return publicKey, nil
    },
    SigningMethod: jwt.SigningMethodRS256,
})


