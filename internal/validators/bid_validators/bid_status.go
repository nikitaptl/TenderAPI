package bid_validators

import (
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/validators"
	"github.com/gorilla/mux"
	"net/http"
)

func ValidateStatusUpdate(r *http.Request) (*models.Bid, error) {
	bid, err := ValidateGettingStatus(r)
	if err != nil {
		return nil, err
	}
	newStatus := r.URL.Query().Get("status")
	if err = ValidateNewBidStatus(bid.Status, newStatus); err != nil {
		return nil, err
	}
	bid.Status = newStatus
	return bid, nil
}

func ValidateGettingStatus(r *http.Request) (*models.Bid, error) {
	vars := mux.Vars(r)
	bidId := vars["bidId"]
	username := r.URL.Query().Get("username")

	bid, err := ValidateBidId(bidId)
	if err != nil {
		return nil, err
	}
	user, err := validators.UserExist(username)
	if err != nil {
		return nil, err
	}
	if _, err = validators.IsUserResponsible(bid.OrganizationID, user.ID); err != nil {
		return nil, err
	}
	return bid, err
}
