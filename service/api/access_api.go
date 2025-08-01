package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Nickname string `json:"nickname"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "invalid request.", http.StatusBadRequest)
		return
	}
	if len(input.Nickname) == 0 {
		CreateJsonError(w, "nickname is requested.", http.StatusBadRequest)
		return
	}
	user, err := rt.db.Login(input.Nickname)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error during the check of the user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if user == "" {
		CreateJsonError(w, "invalid request. Nickname cannot be null. ", http.StatusBadRequest)
		return
	}

	userId, code, err := rt.db.IDFromNICK(user)
	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error during the conversion of the user: "+err.Error(), code)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("token: %d", userId))
	CreateJsonAccessResponse(w, "success", userId, http.StatusOK)
}
