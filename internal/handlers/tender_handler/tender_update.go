package tender_handler

import (
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators"
	"avitoTestTask/internal/validators/tender_validators"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func UpdateTenderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	tenderId := vars["tenderId"]
	username := r.URL.Query().Get("username")

	user, err := validators.UserExist(username)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}
	tender, err := tender_validators.ValidateTenderId(tenderId)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}
	if _, err := validators.IsUserResponsible(tender.OrganizationID, user.ID); err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	updates, err := tender_validators.ValidateUpdates(r.Body)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}
	updates["version"] = tender.Version + 1

	tender, err = updateTender(tender, updates)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	log.Printf("Тендер успешно обновлён: %v", tender)
	json.NewEncoder(w).Encode(tender)
}
