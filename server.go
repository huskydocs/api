package main

import (
	"encoding/json"
	"fmt"
	"github.com/huskydocs/engine/persistence"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello, Huskydocs ^_^\n")
}

func (ph *PersistenceHandler) CreateSubject(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Printf("Received request to create subject: %s \n", p.ByName("subject"))
	decoder := json.NewDecoder(r.Body)

	var s Subject
	err := decoder.Decode(&s)

	if err != nil {
		fmt.Printf("Error decoding JSON payload: %v", err)
		fmt.Fprint(w, "Error decoding JSON payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Decoded requst body with email: %s \n", s.email)
	subject := &persistence.Subject{Username: p.ByName("subject"), Email: s.email}
	err = ph.PS.CreateSubject(subject)
	if err != nil {
		fmt.Printf("Error persisting new subject: %v", err)
		fmt.Fprint(w, "Error persisting new subject")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Successfully created subject")
	w.WriteHeader(http.StatusCreated)
}

func (ph *PersistenceHandler) Projects(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
}

func (ph *PersistenceHandler) Project(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func (ph *PersistenceHandler) CreateProject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func (ph *PersistenceHandler) DeleteProject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func (ph *PersistenceHandler) Documents(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func (ph *PersistenceHandler) Document(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func (ph *PersistenceHandler) CreateDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func (ph *PersistenceHandler) UpdateDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func (ph *PersistenceHandler) DeleteDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

type Subject struct {
	email string
}

type PersistenceHandler struct {
	PS *persistence.PersistenceSession
}

func main() {

	persistenceSession := persistence.Init()
	persistenceHandler := &PersistenceHandler{PS: persistenceSession}

	router := httprouter.New()
	router.GET("/", Index)
	router.PUT("/subject/:subject", persistenceHandler.CreateSubject)

	router.GET("/project/:subject", persistenceHandler.Projects)
	router.GET("/project/:subject/:project", persistenceHandler.Project)
	router.PUT("/project/:subject/:project", persistenceHandler.CreateProject)
	router.DELETE("/project/:subject/:project", persistenceHandler.DeleteProject)

	router.GET("/document/:subject/:project", persistenceHandler.Documents)
	router.GET("/document/:subject/:project/*document", persistenceHandler.Document)
	router.PUT("/document/:subject/:project/*document", persistenceHandler.CreateDocument)
	router.POST("/document/:subject/:project/*document", persistenceHandler.UpdateDocument)
	router.DELETE("/document/:subject/:project/*document", persistenceHandler.DeleteDocument)

	log.Fatal(http.ListenAndServe(":8080", router))
}
