package handler

import (
	"encoding/json"
	"net/http"
	"tn-rest/cmd/server/service"

	_ "github.com/mattn/go-sqlite3"
)

type NewResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data ,omitempty"`
	Error   any    `json:"error ,omitempty"`
}

type NationalParkHandler struct {
	Service *service.NationalParkService
}

func (h NationalParkHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload := service.CreateNationalParkInput{}

	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateNationalPark(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(NewResponse{Message: "successfully create a national park"})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h NationalParkHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Service.GetNationalParks()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Message string                     `json:"message"`
		Data    []service.GetNationalParks `json:"data"`
	}{
		Message: "successfully get national parks",
		Data:    rows,
	}

	res, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
