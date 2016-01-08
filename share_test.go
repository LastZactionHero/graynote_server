package main

import "testing"

func TestCreateShare(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")
	note := createNote(user, "title", "body")
	permissions := "readwrite"

	share := createShare(note, permissions)
	if share.NoteID != note.ID {
		t.Errorf("Expected NoteID to be %q, got %q", note.ID, share.NoteID)
	}
	if share.Permissions != permissions {
		t.Errorf("Expected permissions to eq %q, got %q", permissions, share.Permissions)
	}
	if len(share.AuthKey) == 0 {
		t.Errorf("Expected AuthKey to be present")
	}
}

func TestCreateSharePermissions(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")
	note := createNote(user, "title", "body")

	shareReadWrite := createShare(note, "readwrite")
	if shareReadWrite == nil {
		t.Errorf("Expected shareReadWrite to be created")
	}

	shareRead := createShare(note, "read")
	if shareRead == nil {
		t.Errorf("Expected shareRead to be created")
	}

	shareBadPermission := createShare(note, "garbage")
	if shareBadPermission != nil {
		t.Errorf("Expected shareBadPermission not to be created")
	}
}

func TestFindShareById(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	share := factoryCreateShare("readwrite")

	foundShare := findShareByID(int64(share.ID))
	if foundShare.ID != share.ID {
		t.Errorf("Expected share id %q, got %q", share.ID, foundShare.ID)
	}
}

func TestFindShareByAuthKey(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	share := factoryCreateShare("readwrite")

	foundShare := findShareByAuthKey(share.AuthKey)
	if foundShare.ID != share.ID {
		t.Errorf("Expected share id %q, got %q", share.ID, foundShare.ID)
	}
}

func TestDestroyShare(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	share := factoryCreateShare("readwrite")
	share.Destroy()

	if findShareByID(int64(share.ID)) != nil {
		t.Errorf("Expected share to be deleted")
	}
}

func TestValidateSharePermissions(t *testing.T) {
	if !ValidateSharePermission("readwrite") {
		t.Errorf("Expected readwrite to be valid")
	}
	if !ValidateSharePermission("read") {
		t.Errorf("Expected read to be valid")
	}
	if ValidateSharePermission("garbage") {
		t.Errorf("Expected garbage to be invalid")
	}
}
