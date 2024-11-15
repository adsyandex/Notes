// auth_test.go
package main

import "testing"

func TestHashPassword(t *testing.T) {
    password := "secret"
    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    if !CheckPasswordHash(password, hash) {
        t.Fatalf("Password does not match hash")
    }
}
