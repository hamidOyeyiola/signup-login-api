package controller

import (
	model "github.com/hamidOyeyiola/registration-and-login/models"
)

type Creater interface {
	CreateCreate(q model.SQLQueryToInsert) (h *Header, b *Body, ok bool)
}

type Retriever interface {
	Retrieve(m model.Model) (h *Header, b *Body, ok bool)
}

type Updater interface {
	Update(q model.SQLQueryToUpdate) (h *Header, b *Body, ok bool)
}

type Deleter interface {
	Delete(q model.SQLQueryToDelete) (h *Header, b *Body, ok bool)
}

type CreateRetrieveUpdateDeleter interface {
	Creater
	Retriever
	Updater
	Deleter
}
