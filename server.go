package main

import (
	"encoding/json"
	"fmt"
	"github.com/huskydocs/engine/persistence"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello, Huskydocs ^_^\n")
}

func (ph *PersistenceHandler) CreateSubject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Printf("Received request to create subject: %s \n", params.ByName("subject"))
	s, err := decodeSubject(r)
	if err != nil {
		handleDecodingError(w, err)
		return
	}

	s.Username = params.ByName("subject")

	err = persistSubject(s, ph.PS)
	if err != nil {
		fmt.Printf("Error persisting new subject: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred while creating the subject"))
		return
	}
	w.Write([]byte("Subject created successfully"))
}

func (ph *PersistenceHandler) Projects(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Printf("Received request for all projects owned by: %s \n", params.ByName("subject"))

	owner, err := ph.PS.Subject(params.ByName("subject"))
	if err != nil {
		fmt.Printf("Error retrieving all projects owned by %s: %v", params.ByName("subject"), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred while creating the subject"))
		return
	}

	projects, err := ph.PS.Projects(&owner)
	p := make([]Project, len(projects))

	for i := range projects {
		p[i] = Project{Subject: owner.Username, Name: projects[i].Name, Description: projects[i].Description}
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(p)
}

func (ph *PersistenceHandler) Project(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (ph *PersistenceHandler) CreateProject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Printf("Received request to create Project: %s/%s \n", params.ByName("subject"), params.ByName("project"))
	p, err := decodeProject(r)
	if err != nil {
		handleDecodingError(w, err)
		return
	}

	p.Subject = params.ByName("subject")
	p.Name = params.ByName("project")

	err = persistProject(p, ph.PS)
	if err != nil {
		fmt.Printf("Error persisting new project: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred while creating the project"))
		return
	}
	w.Write([]byte("Project created successfully"))
}

func (ph *PersistenceHandler) DeleteProject(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (ph *PersistenceHandler) Documents(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (ph *PersistenceHandler) Document(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (ph *PersistenceHandler) CreateDocument(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Printf("Received request to create Document: %s/%s/%s \n", params.ByName("subject"), params.ByName("project"), params.ByName("document"))

	d, err := readDocument(r)
	if err != nil {
		handleBodyReadError(w, err)
		return
	}

	d.Subject = params.ByName("subject")
	d.Project = params.ByName("project")
	d.Path = params.ByName("document")

	err = persistDocument(d, ph.PS)
	if err != nil {
		fmt.Printf("Error persisting new document: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred while creating the document"))
		return
	}
	w.Write([]byte("Document created successfully"))
}

func (ph *PersistenceHandler) UpdateDocument(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (ph *PersistenceHandler) DeleteDocument(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

type Subject struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Project struct {
	Subject     string `json:"subject"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Document struct {
	Subject string `json:"subject"`
	Project string `json:"project"`
	Path    string `json:"path"`
	Content []byte
}

func decodeSubject(r *http.Request) (Subject, error) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var s Subject
	err := decoder.Decode(&s)
	return s, err
}

func decodeProject(r *http.Request) (Project, error) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var p Project
	err := decoder.Decode(&p)
	return p, err
}

func readDocument(r *http.Request) (Document, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Document{}, err
	}
	doc := Document{Content: body}
	return doc, nil
}

func persistSubject(s Subject, ps *persistence.PersistenceSession) error {
	subject := &persistence.Subject{Username: s.Username, Email: s.Email}
	err := ps.CreateSubject(subject)
	return err
}

func persistProject(p Project, ps *persistence.PersistenceSession) error {
	owner, err := ps.Subject(p.Subject)
	if err != nil {
		return err
	}

	ownerRef := mgo.DBRef{Collection: "subject", Id: owner.Id}
	project := &persistence.Project{Owner: ownerRef, Name: p.Name, Description: p.Description}
	err = ps.CreateProject(project)
	return err
}

func persistDocument(d Document, ps *persistence.PersistenceSession) error {
	owner, err := ps.Subject(d.Subject)
	if err != nil {
		return err
	}

	project, err := ps.Project(&owner, d.Project)
	if err != nil {
		return err
	}
	projectRef := mgo.DBRef{Collection: "project", Id: project.Id}

	document := &persistence.Document{Project: projectRef, Path: d.Path}
	err = ps.CreateDocument(document)
	return err
}

func handleDecodingError(w http.ResponseWriter, err error) {
	fmt.Printf("Error decoding JSON payload: %v", err)
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Error decoding JSON payload: %v", err)
}

func handleBodyReadError(w http.ResponseWriter, err error) {
	fmt.Printf("Error reading request body: %v", err)
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Error reading request body: %v", err)
}
