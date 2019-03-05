package users

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.
import (
	"testing"
	"golang.org/x/crypto/bcrypt"
)

func TestUserValidate(t *testing.T) {
	cases := []struct {
		name string
		hint string
		nu NewUser
		expectedErr bool
	}{
		{
			"Valid New User",
			"Remember to return nil if everything is valid",
			NewUser{
				"test@example.com",
				"1234567",
				"1234567",
				"test1",
				"tester",
				"testing",
			},
			false,
		},
		{
			"Invalid Email",
			"Remember to use mail.ParseAddress",
			NewUser{
				"invalid email",
				"1234567",
				"1234567",
				"test2",
				"tester",
				"testing",
			},
			true,
		},
		{
			"Short password",
			"remember to check the password length",
			NewUser{
				"test@example.com",
				"123",
				"1234567",
				"test3",
				"tester",
				"testing",
			},
			true,
		},
		{
			"password confirmation mismatch",
			"remember to check if the password and passwordConf are matching",
			NewUser{
				"test@example.com",
				"1234567",
				"1234568",
				"test4",
				"tester",
				"testing",
			},
			true,
		},
		{
			"zero length username",
			"Remember to check if the username is non-zero length",
			NewUser{
				"test@example.com",
				"1234567",
				"1234567",
				"",
				"tester",
				"testing",
			},
			true,
		},
		{
			"username including spaces",
			"Remember to check if the username includes spaces",
			NewUser{
				"test@example.com",
				"1234567",
				"1234567",
				"test  5678",
				"tester",
				"testing",
			},
			true,
		},
	}
	for _, c := range cases {
		err := c.nu.Validate()
		if !c.expectedErr && err != nil {
			t.Errorf("case %s: unexpected error: %v\nHINT: %s", c.name, err, c.hint)
		}
		if c.expectedErr && err == nil {
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
	}
}

func TestNewUser_ToUser(t *testing.T) {
	cases := []struct {
		name string
		hint string
		newUser NewUser
		testNewUser NewUser
		expectedErr bool
		expectedPhotoErr bool
	}{
		{
			"Valid New User",
			"Remember to return nil if everything is valid",
			NewUser{
				"test@example.com",
				"1234567",
				"1234567",
				"test",
				"tester",
				"testing",
			},
			NewUser{
				"TEST@example.com",
				"1234567",
				"1234567",
				"test",
				"tester",
				"testing",
			},
			false,
			false,
		},
		{
			"Invalid Email",
			"Remember to use mail.ParseAddress in validate function",
			NewUser{
				"invalid email",
				"1234567",
				"1234567",
				"test",
				"tester",
				"testing",
			},
			NewUser{
				"TEST@example.com",
				"1234567",
				"1234567",
				"test",
				"tester",
				"testing",
			},
			true,
			true,
		},
		{
			"email with spaces",
			"Remember to use mail.ParseAddress in validate function",
			NewUser{
				" test@example.com ",
				"1234567",
				"1234567",
				"test",
				"tester",
				"testing",
			},
			NewUser{
				"TEST@example.com",
				"1234567",
				"1234567",
				"test",
				"tester",
				"testing",
			},
			false,
			false,
		},
	}
	for _, c := range cases {
		user, err := c.newUser.ToUser()
		if !c.expectedErr && err != nil {
			t.Errorf("case %s: unexpected error: %v\nHINT: %s", c.name, err, c.hint)
		}
		if c.expectedErr && err == nil {
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
		expectedUser ,_ := c.testNewUser.ToUser()
		if !c.expectedPhotoErr && expectedUser.PhotoURL != user.PhotoURL {
			t.Errorf("case %s: unexpected error: expected %s, but got %s", c.name, expectedUser.PhotoURL, user.PhotoURL)
		}
	}
}

func TestUser_FullName(t *testing.T) {
	cases := []struct {
		name string
		hint string
		nu *User
		expectedOutput string
	}{
		{
			"Both first and last name filled",
			"Should be in form of <firstName> <lastName>",
			&User{
				0,
				"test@example.com",
				[]byte("1234567"),
				"test",
				"tester",
				"testing",
				gravatarBasePhotoURL,
			},
			"tester testing",
		},
		{
			"First name only",
			"Should only be first name with no spaces",
			&User{
				0,
				"test@example.com",
				[]byte("1234567"),
				"test",
				"tester",
				"",
				gravatarBasePhotoURL,
			},
			"tester",
		},
		{
			"First name only",
			"Should only be first name with no spaces",
			&User{
				0,
				"test@example.com",
				[]byte("1234567"),
				"test",
				"",
				"testing",
				gravatarBasePhotoURL,
			},
			"testing",
		},
		{
			"First name only",
			"Should only be first name with no spaces",
			&User{
				0,
				"test@example.com",
				[]byte("1234567"),
				"test",
				"",
				"",
				gravatarBasePhotoURL,
			},
			"",
		},
	}
	for _, c := range cases {
		fullName := c.nu.FullName()
		if fullName != c.expectedOutput {
			t.Errorf("case %s: expected to get %s, but got %s\nHINT: %s", c.name, c.expectedOutput, fullName, c.hint)
		}
	}
}

func TestUser_Authenticate(t *testing.T) {
	bCryptPass,_ := bcrypt.GenerateFromPassword([]byte("1234567"), bcryptCost)
	cases := []struct{
		name string
		hint string
		password string
		u User
		expectedErr bool
	}{
		{
			"valid password",
			"look into compareHashAndPassword from bcrypt library",
			"1234567",
			User{
				0,
				"test@example.com",
				bCryptPass,
				"test",
				"tester",
				"testing",
				gravatarBasePhotoURL,
			},
			false,
		},
		{
			"incorrect password",
			"look into compareHashAndPassword from bcrypt library",
			"1234568",
			User{
				0,
				"test@example.com",
				bCryptPass,
				"test",
				"tester",
				"testing",
				gravatarBasePhotoURL,
			},
			true,
		},
		{
			"empty password",
			"look into compareHashAndPassword from bcrypt library",
			"",
			User{
				0,
				"test@example.com",
				bCryptPass,
				"test",
				"tester",
				"testing",
				gravatarBasePhotoURL,
			},
			true,
		},
	}
	for _, c := range cases {
		err := c.u.Authenticate(c.password)
		if !c.expectedErr && err != nil {
			t.Errorf("case %s: unexpected error: %v\nHINT: %s", c.name, err, c.hint)
		}
		if c.expectedErr && err == nil {
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
	}
}

func TestUser_ApplyUpdates(t *testing.T) {
	cases := []struct {
		name string
		hint string
		u User
		update Updates
		expectedErr bool
	}{
		{
			"valid update with both first and last name filled",
			"make sure to check the length of both fields from update object",
			User{
				0,
				"test@example.com",
				[]byte{1,2,3,4},
				"tester",
				"test",
				"testing",
				gravatarBasePhotoURL,
			},
			Updates{
				"testing",
				"test",
			},
			false,
		},
		{
			"valid update with just first name filled",
			"make sure to check the length of both fields from update object",
			User{
				0,
				"test@example.com",
				[]byte{1,2,3,4},
				"tester",
				"test",
				"testing",
				gravatarBasePhotoURL,
			},
			Updates{
				"testing",
				"",
			},
			false,
		},
		{
			"valid update with just last name filled",
			"make sure to check the length of both fields from update object",
			User{
				0,
				"test@example.com",
				[]byte{1,2,3,4},
				"tester",
				"test",
				"testing",
				gravatarBasePhotoURL,
			},
			Updates{
				"",
				"test",
			},
			false,
		},
		{
			"invalid update",
			"make sure to check the length of both fields from update object",
			User{
				0,
				"test@example.com",
				[]byte{1,2,3,4},
				"tester",
				"test",
				"testing",
				gravatarBasePhotoURL,
			},
			Updates{
				"",
				"",
			},
			true,
		},
	}
	for _, c := range cases {
		err := c.u.ApplyUpdates(&c.update)
		if !c.expectedErr && err != nil {
			t.Errorf("case %s: unexpected error: %v\nHINT: %s", c.name, err, c.hint)
		}
		if c.expectedErr && err == nil {
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
	}
}
