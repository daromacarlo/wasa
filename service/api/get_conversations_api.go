package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := r.Header.Get("Authorization")
	exist, err := rt.VerifyToken(token)
	if exist == false {
		CreateJsonError(w, "Error: Token not valid", 401)
		return
	}
	if err != nil {
		CreateJsonError(w, "Error checking token", http.StatusInternalServerError)
		return
	}

	CallingUser := ps.ByName("user")
	list, errorCode, err := rt.db.GetConversationsDB(CallingUser)

	if !errors.Is(err, nil) {
		CreateJsonError(w, "Error while retrieving conversations: "+err.Error(), errorCode)
		return
	}

	if len(list) == 0 {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		CreateJsonError(w, "Error while encoding JSON response: "+err.Error(), http.StatusInternalServerError)
	}
}
