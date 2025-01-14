package bcrypt_test

import (
	"myadmin/util/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	formPassword := "123456"
	savePassword := bcrypt.GeneratePassword(formPassword)
	seccendPassword := bcrypt.GeneratePassword(formPassword)
	if savePassword == seccendPassword {
		t.Log(savePassword, seccendPassword)
		t.Fail()
	}

	if !bcrypt.ComparePassword(formPassword, savePassword) {
		t.Fail()
	}

	if !bcrypt.ComparePassword(formPassword, seccendPassword) {
		t.Fail()
	}

	if bcrypt.ComparePassword(formPassword+"@", savePassword) {
		t.Fail()
	}

}
