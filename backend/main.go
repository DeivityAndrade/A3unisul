package main

import (
	insumos "backend/services/insumo-a3"
	materiaprima "backend/services/materiaprima-a3"
	produto "backend/services/produto-a3"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	//Rotas de produtos
	router.HandleFunc("/produto", produto.Create).Methods(http.MethodPost)
	router.HandleFunc("/produto", produto.FindAll).Methods(http.MethodGet)
	router.HandleFunc("/produto/{id}", produto.Find).Methods(http.MethodGet)
	router.HandleFunc("/produto/{id}", produto.Update).Methods(http.MethodPut)
	router.HandleFunc("/produto/{id}", produto.Delete).Methods(http.MethodDelete)

	//Rotas de Materia Prima
	router.HandleFunc("/materia-prima", materiaprima.Create).Methods(http.MethodPost)
	router.HandleFunc("/materia-prima", materiaprima.FindAll).Methods(http.MethodGet)
	router.HandleFunc("/materia-prima/{id}", materiaprima.Find).Methods(http.MethodGet)
	router.HandleFunc("/materia-prima/{id}", materiaprima.Update).Methods(http.MethodPut)
	router.HandleFunc("/materia-prima/{id}", materiaprima.Delete).Methods(http.MethodDelete)

	//Rotas de Insumos
	router.HandleFunc("/insumos", insumos.Create).Methods(http.MethodPost)
	router.HandleFunc("/insumos", insumos.FindAll).Methods(http.MethodGet)
	router.HandleFunc("/insumos/{id}", insumos.Find).Methods(http.MethodGet)
	router.HandleFunc("/insumos/{id}", insumos.Update).Methods(http.MethodPut)
	router.HandleFunc("/insumos/{id}", insumos.Delete).Methods(http.MethodDelete)

	fmt.Println("Escutando na porta 5011")
	log.Fatal(http.ListenAndServe(":5020", router))
}
