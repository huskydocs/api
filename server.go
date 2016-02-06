package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "log"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Hello, Huskydocs ^_^\n")
}

func Projects(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func Project(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func CreateProject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func DeleteProject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func Documents(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func Document(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func CreateDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func UpdateDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func DeleteDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func main() {
    router := httprouter.New()
    router.GET("/", Index)

    router.GET("/project/:subject", Projects)
    router.GET("/project/:subject/:project", Project)
    router.PUT("/project/:subject/:project", CreateProject)
    router.DELETE("/project/:subject/:project", DeleteProject)
    
    router.GET("/document/:subject/:project", Documents)
    router.GET("/document/:subject/:project/*document", Document)
    router.PUT("/document/:subject/:project/*document", CreateDocument)
    router.POST("/document/:subject/:project/*document", UpdateDocument)
    router.DELETE("/document/:subject/:project/*document", DeleteDocument)
    
    log.Fatal(http.ListenAndServe(":8080", router))
}
