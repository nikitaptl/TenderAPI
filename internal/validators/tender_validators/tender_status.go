package tender_validators

import (
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/validators"
	"github.com/gorilla/mux"
	"net/http"
)

func ValidateStatusUpdate(r *http.Request) (*models.Tender, error) {
	tender, err := ValidateGettingStatus(r)
	if err != nil {
		return nil, err
	}

	status := r.URL.Query().Get("status")
	if err = ValidateNewTenderStatus(tender.Status, status); err != nil {
		return nil, err
	}
	tender.Status = status
	return tender, nil
}

func ValidateGettingStatus(r *http.Request) (*models.Tender, error) {
	vars := mux.Vars(r)
	tenderId := vars["tenderId"]
	username := r.URL.Query().Get("username")

	tender, err := ValidateTenderId(tenderId)
	if err != nil {
		return nil, err
	}
	user, err := validators.UserExist(username)
	if err != nil {
		return nil, err
	}
	if _, err = validators.IsUserResponsible(tender.OrganizationID, user.ID); err != nil {
		return nil, err
	}
	return tender, nil
}
