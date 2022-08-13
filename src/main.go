package main

import (
	"context"
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"

	"github.com/go-session/session"
	"github.com/gorilla/mux"

	model "github.com/hamidOyeyiola/registration-and-login/models"
	"github.com/hamidOyeyiola/registration-and-login/utils"

	controller "github.com/hamidOyeyiola/registration-and-login/controllers"
)

const (
	dataSource = "hamid:@tcp(localhost:3306)/registrationandlogin"
)

func main() {
	r := mux.NewRouter()

	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost",
	})
	manager.MapClientStorage(clientStore)
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetInternalErrorHandler(func(err error) (er *errors.Response) {
		log.Println("Internal Error: ", err.Error())
		return
	})

	srv.SetUserAuthorizationHandler(userAuthorization)
	srv.SetResponseErrorHandler(func(er *errors.Response) {
		log.Println("Internal Error: ", er.Error.Error())
		return
	})

	r.HandleFunc("/authorize", func(rw http.ResponseWriter, req *http.Request) {
		err := srv.HandleAuthorizeRequest(rw, req)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	})

	r.HandleFunc("/token", func(rw http.ResponseWriter, req *http.Request) {
		srv.HandleTokenRequest(rw, req)
	})

	r.HandleFunc("/", index)
	r.HandleFunc("/login", login)

	go func() {
		err := http.ListenAndServe("localhost:9096", r)
		if err != nil {
			log.Println(err)
		}
	}()

	fmt.Println("Server is now ready to take your orders!")
	e := make(chan os.Signal, 1)
	signal.Notify(e, os.Interrupt)

	<-e

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting Down!")
	os.Exit(1)
}

func userAuthorization(rw http.ResponseWriter, req *http.Request) (userID string, err error) {
	store, err := session.Start(nil, rw, req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	uid, ok := store.Get("userID")
	if !ok {
		if req.Form == nil {
			req.ParseForm()
		}
		store.Set("ReturnUri", req.Form)
		store.Save()
		rw.Header().Set("Location", "/login")
		rw.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	store.Delete("userID")
	store.Save()
	return
}

func index(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		fname := req.FormValue("fname")
		lname := req.FormValue("lname")
		phone := req.FormValue("phone")
		email := req.FormValue("email")
		password := req.FormValue("password")
		passwordConfirm := req.FormValue("passwordConfirm")
		if password != passwordConfirm {
			return
		}
		user := model.User{
			FirstName: fname,
			LastName:  lname,
			PhoneNo:   phone,
			Email:     utils.EmailAddress(email),
			Password:  password,
		}
		q, ok := user.Insert()
		if !ok {

		}
		cc := controller.NewMySQLController(dataSource)
		response := new(controller.Response)
		h, b, ok := cc.Create(q)
		if ok {
			o, _ := json.MarshalIndent(user, "", " ")
			b = new(controller.Body).
				AddContentType("application/json").
				AddContent(string(o))
		}
		response.AddHeader(h).
			AddBody(b).
			Write(rw)
		return
	}
	html(rw, req, "html/index.html")
}

func login(rw http.ResponseWriter, req *http.Request) {
	store, err := session.Start(nil, rw, req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {
		email := req.FormValue("email")
		password := req.FormValue("password")
		user := model.User{
			Email:    utils.EmailAddress(email),
			Password: password,
		}

		cc := controller.NewMySQLController(dataSource)
		cc.Retrieve(&user)
		ok := utils.EncryptPassword(password) == user.Password
		if ok {
			controller.GetStatusBadRequestRes()
			return
		}
		var form url.Values
		if v, ok := store.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		u := new(url.URL)
		u.Path = "/authorize"
		u.RawQuery = form.Encode()
		rw.Header().Set("Location", u.String())
		rw.WriteHeader(http.StatusFound)
		store.Delete("Form")
		store.Set("userID", string(user.Email))
		return
	}
	html(rw, req, "html/login.html")
}

func html(rw http.ResponseWriter, req *http.Request, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	st, _ := file.Stat()
	http.ServeContent(rw, req, file.Name(), st.ModTime(), file)
}
