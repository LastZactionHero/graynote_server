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

func TestCreateUser(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	email := "user@site.com"
	password := "mypassword"
	userForm := UserRegisterForm{Email: email, Password: password}
	user := createUser(&userForm)

	if user.Email != email {
		t.Errorf("Expected email to be %q, got %q", email, user.Email)
	}
	if len(user.PasswordHash) == 0 {
		t.Errorf("Expected passwordHash to be present")
	}
}

func TestFindUserByID(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")

	foundUser := findUserByID(int64(user.ID))
	if foundUser.ID != user.ID {
		t.Errorf("Expected user %q, got  %q", user, foundUser)
	}
}

func TestFindUserByEmail(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	email := "user@site.com"
	user := factoryCreateUser(email)

	foundUser := findUserByEmail(email)
	if foundUser.ID != user.ID {
		t.Errorf("Expected user %q, got  %q", user, foundUser)
	}
}

func TestFindUserByAuthToken(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	email := "user@site.com"
	user := factoryCreateUser(email)

	foundUser := findUserByAuthToken(user.AuthToken)
	if foundUser.ID != user.ID {
		t.Errorf("Expected user %q, got  %q", user, foundUser)
	}
}

func TestValidPasswordForUser(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	email := "user@site.com"
	user := factoryCreateUser(email)

	if !user.validPasswordForUser("password") {
		t.Errorf("Expected password to be valid for user")
	}
}
