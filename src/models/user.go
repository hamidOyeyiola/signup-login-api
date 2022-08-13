package model

import (
	"fmt"

	"github.com/hamidOyeyiola/registration-and-login/interfaces"
	"github.com/hamidOyeyiola/registration-and-login/utils"
)

type User struct {
	FirstName string             `json:"firstname"`
	LastName  string             `json:"lastname"`
	Email     utils.EmailAddress `json:"email"`
	PhoneNo   string             `json:"phoneno"`
	Password  string             `json:"password"`
	CreatedOn string             `json:"-"`
	UpdatedOn string             `json:"-"`
	ID        int                `json:"-"`
}

func (u *User) Insert() (ins SQLQueryToInsert, ok bool) {
	if ok := u.Email.IsValid() && len(u.Password) >= 8; ok {
		ins = SQLQueryToInsert(fmt.Sprintf("INSERT INTO users(firstname, lastname, email, phoneno, password, createdOn,updatedOn,id) VALUES ('%s','%s','%s','%s','%s','%s','%s',%d)",
			u.FirstName, u.LastName, u.Email, u.PhoneNo, utils.EncryptPassword(u.Password), utils.NewDate(), "", 0))
	}
	return
}

func (u *User) Update() (upt SQLQueryToUpdate, ok bool) {
	if ok = u.Email.IsValid() && len(u.Password) >= 8; !ok {
		return
	}
	upt.ID = string(u.Email)
	if u.FirstName != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET firstname=? WHERE email=?"))
		upt.Values = append(upt.Values, u.FirstName)
	}
	if u.LastName != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET lastname=? WHERE email=?"))
		upt.Values = append(upt.Values, u.LastName)
	}
	if u.PhoneNo != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET phoneno=? WHERE email=?"))
		upt.Values = append(upt.Values, u.PhoneNo)
	}
	if u.Password != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET password=? WHERE email=?"))
		upt.Values = append(upt.Values, utils.EncryptPassword(u.Password))
	}
	upt.Stmts = append(upt.Stmts, string("UPDATE users SET updatedOn=? WHERE email=?"))
	date := fmt.Sprintf("%s", utils.NewDate())
	upt.Values = append(upt.Values, date)
	return
}

func (u *User) FromQueryResult(s interfaces.Iterator) error {
	s.Next()
	err := s.Scan(&u.FirstName, &u.LastName, &u.Email, &u.PhoneNo, &u.Password, &u.CreatedOn, &u.UpdatedOn, &u.ID)
	return err
}

func (u *User) Select() (sel SQLQueryToSelect, ok bool) {
	sel, ok = SQLQueryToSelect(fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", u.Email)), true
	return
}

func (u *User) Delete() (del SQLQueryToDelete, ok bool) {
	del, ok = SQLQueryToDelete(fmt.Sprintf("DELETE FROM users WHERE email = '%s'", u.Email)), true
	return
}
