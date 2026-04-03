package handler

import "net/http"

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Products"))
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Product"))
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create Product"))
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update Product"))
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete Product"))
}
