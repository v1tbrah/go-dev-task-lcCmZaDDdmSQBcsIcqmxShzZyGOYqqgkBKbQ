package api

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"

	dberr "go-dev-task-lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ/internal/phone/storage/error"
)

func (a *API) getPhone(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is missing", http.StatusBadRequest)
		return
	}

	idForGetPhone, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	phone, err := a.storage.GetPhone(idForGetPhone)
	if err != nil {
		if errors.Is(err, dberr.ErrPhoneIsNotFound) {
			http.Error(w, "phone is not found", http.StatusNotFound)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	if unlucky() {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func unlucky() bool {
	return rand.Intn(2) == 0
}
