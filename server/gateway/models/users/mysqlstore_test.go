package users

import (
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"fmt"
)

//successful test
func TestMySQLStore_Insert_GetByID_Delete(t *testing.T) {
	//create a new sql mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	//ensure it's closed at the end of the test
	defer db.Close()

	//construct a new task and a set of rows to return
	expectedUser := &User{
		ID: 1,
		Email: "test@example.com",
		PassHash: []byte("1010"),
		UserName: "tester",
		FirstName: "test",
		LastName: "testing",
		PhotoURL: gravatarBasePhotoURL,
	}

	//tell sqlmock that we expect the function to execute a
	//a particular SQL query
	expectedSQL := "insert into users(email, username, passHash, firstName, lastName, photoUrl) values(?,?,?,?,?,?)"
	//mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(expectedUser.Email, expectedUser.UserName,
		expectedUser.PassHash, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL).
			WillReturnResult(sqlmock.NewResult(1, 1))
	//mock.ExpectCommit()
	//construct a new MySQLStore using the mock db
	store := NewMySQLStore(db)

	//now execute insert method
	user, err := store.Insert(expectedUser)
	if err != nil {
		t.Errorf("unexpected error inserting a user, but got %v\n", err)
	}

	if user == nil {
		t.Errorf("expected to get a user, but got none")
	}

	if user.ID != 1 {
		t.Errorf("expected the user ID to be 1, but got %d", user.ID)
	}

	rows := sqlmock.NewRows([]string{"id", "email", "username", "passHash", "firstName", "lastName", "photoUrl"})
	rows.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.UserName,
		expectedUser.PassHash, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)
	expectedSQL = "select * from users where id=?"
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(string(expectedUser.ID)).WillReturnRows(rows)
	mock.ExpectCommit()
	user, err = store.GetByID(expectedUser.ID)
	if err != nil {
		t.Errorf("unexpected error getting a user by ID, but got %v\n", err)
	}

	if user == nil {
		t.Errorf("expected a user by id, but got none")
	}

	expectedSQL = "DELETE FROM users WHERE id=?"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(expectedUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.Delete(expectedUser.ID)
	if err != nil {
		t.Errorf("unexpected error deleting a user, but got %v\n", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations for deletion: %v", err)
	}
}

//test for transaction errors
func TestMySQLStore_Transaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	//construct a new task and a set of rows to return
	expectedUser := &User{
		ID: 1,
		Email: "test@example.com",
		PassHash: []byte("1010"),
		UserName: "tester",
		FirstName: "test",
		LastName: "testing",
		PhotoURL: gravatarBasePhotoURL,
	}

	store := NewMySQLStore(db)
	_, err = store.Insert(expectedUser)
	if err == nil {
		t.Errorf("expected an error for starting transaction, but got none")
	}

	mock.ExpectBegin()
	_, err = store.Insert(expectedUser)
	if err == nil {
		t.Errorf("expected an error for exec, but got none")
	}

	update := &Updates{
		FirstName: "testing",
		LastName: "test",
	}

	_, err = store.Update(expectedUser.ID, update)
	if err == nil {
		t.Errorf("expected an error for starting transaction, but got none")
	}

	mock.ExpectBegin()
	_, err = store.Update(expectedUser.ID, update)
	if err == nil {
		t.Errorf("expected an error for exec, but got none")
	}

	err = store.Delete(expectedUser.ID)
	if err == nil {
		t.Errorf("expected an error for starting transaction, but got none")
	}

	mock.ExpectBegin()
	err = store.Delete(expectedUser.ID)
	if err == nil {
		t.Errorf("expected an error for exec, but got none")
	}

	_, err = store.GetByID(expectedUser.ID)
	if err == nil {
		t.Errorf("expected an error for starting transaction, but got none")
	}

	_, err = store.GetByEmail(expectedUser.Email)
	if err == nil {
		t.Errorf("expected an error for starting transaction, but got none")
	}

	_, err = store.GetByUserName(expectedUser.UserName)
	if err == nil {
		t.Errorf("expected an error for starting transaction, but got none")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations for deletion: %v", err)
	}
}

func TestMySQLStore_Update_GetById(t *testing.T) {
	//create a new sql mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	//ensure it's closed at the end of the test
	defer db.Close()

	//construct a new task and a set of rows to return
	expectedUser := &User{
		ID: 1,
		Email: "test@example.com",
		PassHash: []byte("1010"),
		UserName: "tester",
		FirstName: "test",
		LastName: "testing",
		PhotoURL: gravatarBasePhotoURL,
	}

	//tell sqlmock that we expect the function to execute a
	//a particular SQL query
	//construct a new MySQLStore using the mock db
	store := NewMySQLStore(db)

	//now execute insert method
	mock.ExpectBegin()
	user, _ := store.Insert(expectedUser)

	expectedSQL := "UPDATE users SET firstName=?, lastName=? where id=?"
	update := &Updates{
		FirstName: "testing",
		LastName: "test",
	}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(update.FirstName,
		update.LastName, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	//check for getByID
	rows := sqlmock.NewRows([]string{"id", "email", "username", "passHash", "firstName", "lastName", "photoUrl"})
	rows.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.UserName,
		expectedUser.PassHash, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)
	expectedSQL = "select * from users where id=?"
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(string(expectedUser.ID)).WillReturnRows(rows)
	mock.ExpectCommit()
	user, err = store.Update(1, update)
	if err != nil {
		t.Errorf("unexpected error updating a user, but got %v\n", err)
	}

	if user == nil {
		t.Errorf("expected a user after update, but got none")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations for deletion: %v", err)
	}
}

func TestMySQLStore_Insert_GetByEmail_DELETE(t *testing.T) {
	//create a new sql mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	//ensure it's closed at the end of the test
	defer db.Close()

	//construct a new task and a set of rows to return
	expectedUser := &User{
		ID: 1,
		Email: "test@example.com",
		PassHash: []byte("1010"),
		UserName: "tester",
		FirstName: "test",
		LastName: "testing",
		PhotoURL: gravatarBasePhotoURL,
	}

	//tell sqlmock that we expect the function to execute a
	//a particular SQL query
	expectedSQL := "insert into users(email, username, passHash, firstName, lastName, photoUrl) values(?,?,?,?,?,?)"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(expectedUser.Email, expectedUser.UserName,
		expectedUser.PassHash, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	//construct a new MySQLStore using the mock db
	store := NewMySQLStore(db)

	//now execute insert method
	user, err := store.Insert(expectedUser)
	if err != nil {
		t.Errorf("unexpected error inserting a user, but got %v\n", err)
	}

	if user == nil {
		t.Errorf("expected to get a user, but got none")
	}

	if user.ID != 1 {
		t.Errorf("expected the user ID to be 1, but got %d", user.ID)
	}

	rows := sqlmock.NewRows([]string{"id", "email", "username", "passHash", "firstName", "lastName", "photoUrl"})
	rows.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.UserName,
		expectedUser.PassHash, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)
	expectedSQL = "select * from users where email=?"
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(expectedUser.Email).WillReturnRows(rows)
	mock.ExpectCommit()
	user, err = store.GetByEmail(expectedUser.Email)
	if err != nil {
		t.Errorf("unexpected error getting a user by ID, but got %v\n", err)
	}

	if user == nil {
		t.Errorf("expected a user by id, but got none")
	}

	expectedSQL = "DELETE FROM users WHERE id=?"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(expectedUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.Delete(expectedUser.ID)
	if err != nil {
		t.Errorf("unexpected error deleting a user, but got %v\n", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations for deletion: %v", err)
	}
}

func TestMySQLStore_Insert_GetByUsername_Delete(t *testing.T) {
	//create a new sql mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	//ensure it's closed at the end of the test
	defer db.Close()

	//construct a new task and a set of rows to return
	expectedUser := &User{
		ID: 1,
		Email: "test@example.com",
		PassHash: []byte("1010"),
		UserName: "tester",
		FirstName: "test",
		LastName: "testing",
		PhotoURL: gravatarBasePhotoURL,
	}

	//tell sqlmock that we expect the function to execute a
	//a particular SQL query
	expectedSQL := "insert into users(email, username, passHash, firstName, lastName, photoUrl) values(?,?,?,?,?,?)"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(expectedUser.Email, expectedUser.UserName,
		expectedUser.PassHash, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	//construct a new MySQLStore using the mock db
	store := NewMySQLStore(db)

	//now execute insert method
	user, err := store.Insert(expectedUser)
	if err != nil {
		t.Errorf("unexpected error inserting a user, but got %v\n", err)
	}

	if user == nil {
		t.Errorf("expected to get a user, but got none")
	}

	if user.ID != 1 {
		t.Errorf("expected the user ID to be 1, but got %d", user.ID)
	}

	rows := sqlmock.NewRows([]string{"id", "email", "username", "passHash", "firstName", "lastName", "photoUrl"})
	rows.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.UserName,
		expectedUser.PassHash, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)
	//test the get by username
	expectedSQL = "select * from users where username=?"
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(expectedUser.UserName).WillReturnRows(rows)
	mock.ExpectCommit()
	user, err = store.GetByUserName(expectedUser.UserName)
	if err != nil {
		t.Errorf("unexpected error getting a user by ID, but got %v\n", err)
	}

	if user == nil {
		t.Errorf("expected a user by id, but got none")
	}

	expectedSQL = "DELETE FROM users WHERE id=?"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(expectedUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.Delete(expectedUser.ID)
	if err != nil {
		t.Errorf("unexpected error deleting a user, but got %v\n", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations for deletion: %v", err)
	}
}

func TestMySQLStore_InsertRollbackOnFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	expectedUser := &User{
		ID: 1,
		Email: "test@example.com",
		PassHash: []byte("1010"),
		UserName: "tester",
		FirstName: "test",
		LastName: "testing",
		PhotoURL: gravatarBasePhotoURL,
	}
	expectedSQL := "insert into users(email, username, passHash, firstName, lastName, photoUrl) values(?,?,?,?,?,?)"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs().WillReturnError(fmt.Errorf("error inserting a new user"))
	mock.ExpectRollback()

	store := NewMySQLStore(db)
	user, err := store.Insert(expectedUser)
	if err == nil {
		t.Errorf("expected error inserting a user, but got none")
	}

	if user != nil {
		t.Errorf("expected no user returned after insertion, but got %v\n", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
