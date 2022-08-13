package model

import (
	interfaces "github.com/hamidOyeyiola/registration-and-login/interfaces"
)

type SQLQueryToInsert string

type SQLQueryToSelect string

type SQLQueryToDelete string

type SQLQueryToUpdate struct {
	Stmts  []string
	Values []string
	ID     string
}

type JSONObject string

type Model interface {
	FromQueryResult(interfaces.Iterator) error
	Insert() (SQLQueryToInsert, bool)
	Update() (SQLQueryToUpdate, bool)
	Select() (SQLQueryToSelect, bool)
	Delete() (SQLQueryToDelete, bool)
}
