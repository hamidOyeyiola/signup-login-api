package controller

import (
	"database/sql"
	"sync"
	"time"

	//"github.com/go-session/session"

	model "github.com/hamidOyeyiola/registration-and-login/models"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLController struct {
	dataSource string
	db         *sql.DB
	conns      int
	mu         sync.Mutex
}

func (sc *MySQLController) Open() bool {
	sc.mu.Lock()
	if sc.db == nil {
		db, err := sql.Open("mysql", sc.dataSource)
		if err != nil {
			return false
		}
		db.SetConnMaxLifetime(3 * time.Minute)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		sc.db = db
	}
	sc.conns++
	sc.mu.Unlock()
	return true
}

func (sc *MySQLController) Close() bool {
	sc.mu.Lock()
	sc.conns--
	if sc.conns == 0 {
		err := sc.db.Close()
		if err != nil {
			return false
		}
		sc.db = nil
	}
	sc.mu.Unlock()
	return true
}

func NewMySQLController(datasrc string) *MySQLController {
	return &MySQLController{dataSource: datasrc}
}

func (sc *MySQLController) Create(q model.SQLQueryToInsert) (h *Header, b *Body, ok bool) {

	ok = sc.Open()
	if !ok {
		h, b = GetStatusFailedDependencyRes()
		ok = false
		return
	}
	defer sc.Close()
	_, err := sc.db.Exec(string(q))
	if err != nil {
		h, b = GetStatusNotModifiedRes()
		ok = false
		return
	}
	h, b = GetStatusCreatedRes()
	return
}

func (sc *MySQLController) Retrieve(m model.Model) (h *Header, b *Body, ok bool) {
	ok = sc.Open()
	if !ok {
		h, b = GetStatusFailedDependencyRes()
		ok = false
		return
	}
	defer sc.Close()
	q, _ := m.Select()
	res, err := sc.db.Query(string(q))
	if err != nil {
		h, b = GetStatusFailedDependencyRes()
		ok = false
		return
	}
	defer res.Close()
	err = m.FromQueryResult(res)
	if err != nil {
		h, b = GetStatusNotFoundRes()
		ok = false
		return
	}
	h, b = GetStatusOKRes()
	return
}

func (sc *MySQLController) Update(q model.SQLQueryToUpdate) (h *Header, b *Body, ok bool) {
	ok = sc.Open()
	if !ok {
		h, b = GetStatusFailedDependencyRes()
		return
	}
	defer sc.Close()
	for j := 0; j < len(q.Stmts); j++ {
		stmt, err := sc.db.Prepare(string(q.Stmts[j]))
		if err != nil {
			h, b = GetStatusNotModifiedRes()
			ok = false
			return
		}
		_, err = stmt.Exec(q.Values[j], q.ID)
		if err != nil {
			h, b = GetStatusNotModifiedRes()
			ok = false
			return
		}
	}
	h, b = GetStatusOKRes()
	return
}

func (sc *MySQLController) Delete(q model.SQLQueryToDelete) (h *Header, b *Body, ok bool) {
	ok = sc.Open()
	if !ok {
		h, b = GetStatusFailedDependencyRes()
		return
	}
	defer sc.Close()
	_, err := sc.db.Exec(string(q))
	if err != nil {
		h, b = GetStatusBadRequestRes()
		ok = false
		return
	}
	h, b = GetStatusOKRes()
	return
}
