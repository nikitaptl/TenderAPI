package tender_handler

import (
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators/tender_validators"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RollbackTender(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	tenderIDStr := vars["tenderId"]
	versionStr := vars["version"]
	username := r.URL.Query().Get("username")

	tender, err := tender_validators.ValidateTenderRollback(tenderIDStr, versionStr, username)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	if err = saveTenderVersion(tender); err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	log.Printf("Тендер успешно откачен до версии %s", versionStr)
	json.NewEncoder(w).Encode(tender)
}
