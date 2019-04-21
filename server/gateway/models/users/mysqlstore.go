package users

import (
	"database/sql"
	"fmt"
	"errors"
	"strings"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/indexes"
)

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore(db *sql.DB) MySQLStore {
	msStore := MySQLStore{
		db,
	}
	return msStore
}

func ParseRows(row *sql.Row) (*User, error) {
	var user User
	err := row.Scan(&user.UserID, &user.Email, &user.UserName, &user.PassHash,
		&user.FirstName, &user.LastName)
	if err != nil {
		if err != sql.ErrNoRows{
			return nil, errors.New(fmt.Sprintf("error scanning a row: %v\n", err))
		}
		return nil, nil
	}
	return &user, nil
}

//helper function that will iterate through a collection of keys to add into trie
func insertIntoTrie(user User, trie indexes.Trie) {
	names := []string{user.FirstName, user.LastName, user.UserName}
	for _, name := range names {
		for _, key := range strings.Split(name, " ") {
			trie.Add(key, user.UserID)
		}
	}
}

//will be used to load all users into trie
func (s *MySQLStore) LoadTrie(trie *indexes.Trie) error {
	rows, err := s.db.Query("select * from users")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		if err = rows.Scan(&user.UserID, &user.UserName, &user.FirstName, &user.LastName); err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return err
		}
		insertIntoTrie(user, *trie)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func (s *MySQLStore) GetByID(id int64) (*User, error) {
	row := s.db.QueryRow("select * from users where userID=?", id)
	user, err := ParseRows(row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *MySQLStore) GetByEmail(email string) (*User, error) {
	row := s.db.QueryRow("select * from users where email=?", email)
	user, err := ParseRows(row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *MySQLStore) GetByUserName(username string) (*User, error) {
	row := s.db.QueryRow("select * from users where username=?", username)
	user, err := ParseRows(row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *MySQLStore) Insert(user *User) (*User, error) {
	insertQuery := "insert into users(email, username, passHash, firstName, lastName, photoUrl) values(?,?,?,?,?,?)"
	res, err := s.db.Exec(insertQuery,user.Email, user.UserName, user.PassHash, user.FirstName, user.LastName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error inserting a new user: %v\n", err))
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting new last inserted id of user: %v\n", err))
	}
	user.UserID = id
	return user, nil
}

func (s *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	updateStmt := "UPDATE users SET firstName=?, lastName=? where userID=?"
	res, err := s.db.Exec(updateStmt, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in updating a user: %v\n", err))
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrUserNotFound
	}
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *MySQLStore) Delete(id int64) error {
	delStmt := "DELETE FROM users WHERE userID=?"
	res, err := s.db.Exec(delStmt, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrUserNotFound
	}
	return nil
}
