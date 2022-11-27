package api

import (
	"math/rand"
	"net/http"
	"time"

	"go-dev-task-lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ/internal/phone/storage"
)

type API struct {
	server  *http.Server
	storage Storage
}

func New() *API {

	newAPI := &API{}

	server := newAPI.newServer()
	newAPI.server = server

	newStorage := storage.New()
	newAPI.storage = newStorage

	return newAPI

}

func (a *API) newServer() *http.Server {

	newServer := &http.Server{Addr: ":3333"}

	router := a.newRouter()

	newServer.Handler = router

	rand.Seed(time.Now().UnixNano())

	return newServer

}

func (a *API) newRouter() http.Handler {

	newRouter := http.NewServeMux()

	newRouter.HandleFunc("/", a.getPhone)

	return newRouter

}

func (a *API) Run() {

	defer a.server.Close()
	a.server.ListenAndServe()

}
