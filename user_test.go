package main

import "testing"

func TestPasswordToHash(t *testing.T) {
	password := "mypassword123"
	expected := "9c87baa223f464954940f859bcf2e233"
	hashedPassword := passwordToHash(password)
	if hashedPassword != expected {
		t.Error("Expected ", expected, " got ", hashedPassword)
	}
}
