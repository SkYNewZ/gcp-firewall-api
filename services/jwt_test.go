package services

import (
	"encoding/json"
	"fmt"
	"testing"
)

type test struct {
	Test     string
	Title    string
	Expected interface{}
}

func TestJWTStruct(t *testing.T) {
	e := JWT{Email: "dummy_mail", EmailVerified: true, Iss: "dummy_iss"}
	jwt, err := json.Marshal(e)
	if err != nil {
		t.Errorf("Unexpected error while decoding token: %s", err)
	}

	expected := `{"iss":"dummy_iss","email":"dummy_mail","email_verified":true}`
	if string(jwt) != expected {
		t.Errorf("JSON marshal is not expected. Got %s want %s", string(jwt), expected)
	}
}

func TestGetUserEmailFromJWT(t *testing.T) {
	tests := []test{
		test{
			Title: "Malformed JWT",
			Test:  "Bearer a",
		},
		test{
			Title: "Malformed JWT",
			Test:  "a.a",
		},
		test{
			Title: "Malformed JWT",
			Test:  "a.a.a.a",
		},
		test{
			// ISS is equal to https://facebook.com
			Test:  "Bearer dummy.ewogICJpc3MiOiAiaHR0cHM6Ly9mYWNlYm9vay5jb20iLAogICJhenAiOiAiMzI1NTU5NDA1NTkuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLAogICJhdWQiOiAiMzI1NTU5NDA1NTkuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLAogICJzdWIiOiAiMTA4NjEwMTM4NjcwMTM2MzEzNTM5IiwKICAiaGQiOiAiZXh0LmFkZW8uY29tIiwKICAiZW1haWwiOiAicXVlbnRpbi5sZW1haXJlQGV4dC5hZGVvLmNvbSIsCiAgImVtYWlsX3ZlcmlmaWVkIjogdHJ1ZSwKICAiYXRfaGFzaCI6ICJjaWVCY3EzWUxwQjdSMVJjc2dCOXV3IiwKICAiaWF0IjogMTU4NTY2NDQyNywKICAiZXhwIjogMTU4NTY2ODAyNwp9.dummy",
			Title: "Invalid issuer",
		},
		test{
			// Email is not verified
			Test:  "Bearer dummy.ewogICJpc3MiOiAiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tIiwKICAiYXpwIjogIjMyNTU1OTQwNTU5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwKICAiYXVkIjogIjMyNTU1OTQwNTU5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwKICAic3ViIjogIjEwODYxMDEzODY3MDEzNjMxMzUzOSIsCiAgImhkIjogImV4dC5hZGVvLmNvbSIsCiAgImVtYWlsIjogInF1ZW50aW4ubGVtYWlyZUBleHQuYWRlby5jb20iLAogICJlbWFpbF92ZXJpZmllZCI6IGZhbHNlLAogICJhdF9oYXNoIjogImNpZUJjcTNZTHBCN1IxUmNzZ0I5dXciLAogICJpYXQiOiAxNTg1NjY0NDI3LAogICJleHAiOiAxNTg1NjY4MDI3Cn0=.dummy",
			Title: "Email not veried",
		},
		test{
			// Valid token
			Test:  "Bearer dummy.ewogICJpc3MiOiAiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tIiwKICAiYXpwIjogIjMyNTU1OTQwNTU5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwKICAiYXVkIjogIjMyNTU1OTQwNTU5LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwKICAic3ViIjogIjEwODYxMDEzODY3MDEzNjMxMzUzOSIsCiAgImhkIjogImV4dC5hZGVvLmNvbSIsCiAgImVtYWlsIjogImR5bW15QGV4dC5hZGVvLmNvbSIsCiAgImVtYWlsX3ZlcmlmaWVkIjogdHJ1ZSwKICAiYXRfaGFzaCI6ICJjaWVCY3EzWUxwQjdSMVJjc2dCOXV3IiwKICAiaWF0IjogMTU4NTY2NDQyNywKICAiZXhwIjogMTU4NTY2ODAyNwp9.dummy",
			Title: "Valid token",
		},
	}

	var expected string
	for _, test := range tests {
		t.Run(test.Title, func(t *testing.T) {
			result, err := GetUserEmailFromJWT(test.Test)

			// Cast okay, we have an error
			if err != nil {
				expected = fmt.Sprintf(`{"code":400,"message":"Invalid Bearer token. Please make sure you are using 'gcloud auth print-identity-token': %s"}`, test.Title)
				if err.Error() != expected {
					t.Errorf("Unexpected error. Got '%s' want '%s'", err.Error(), expected)
				}
			}

			// Cast ko, we don't have error
			if err == nil {
				expected = "dymmy@ext.adeo.com"
				if result != expected {
					t.Errorf("Not wanted email. Got %s, want %s", result, expected)
				}
			}

		})
	}
}

func TestPadBase64Input(t *testing.T) {
	tests := []test{
		test{
			Title:    "Result should be the same",
			Test:     "abcd",
			Expected: "abcd",
		},
		test{
			Title:    "Result should have =",
			Test:     "abc",
			Expected: "abc=",
		},
		test{
			Title:    "Result should have ==",
			Test:     "ab",
			Expected: "ab==",
		},
	}

	for _, test := range tests {
		t.Run(test.Test, func(t *testing.T) {
			result := padBase64Input(test.Test)
			if result != test.Expected {
				t.Errorf("Expected '%s', got '%s'", test.Expected, result)
			}
		})
	}
}
