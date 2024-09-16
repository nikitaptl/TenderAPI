package bid_validators

import (
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func ValidateBidUpdates(r *http.Request) (map[string]interface{}, error) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	bidId := vars["bidId"]
	username := r.URL.Query().Get("username")

	bid, err := ValidateBidAndUser(bidId, username)
	if err != nil {
		return nil, err
	}
	fmt.Printf("bid: %v\n", bid)

	var request models.UpdateBidRequest
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, utils.NewErrorResponse(err.Error(), http.StatusBadRequest)
	}

	updates := make(map[string]interface{})
	if request.Name != nil {
		updates["name"] = *request.Name
	}
	if request.Description != nil {
		updates["description"] = *request.Description
	}
	if len(updates) == 0 {
		return nil, utils.NewErrorResponse("Не переданы поля, которые могут быть обновлены", http.StatusBadRequest)
	}
	// Поскольку предложение обновляется, все предыдущие соглашения становятся недействительными
	quorum, err := FindQuorum(bid.OrganizationTenderID)
	if err != nil {
		return nil, err
	}
	updates["remaining_approvals"] = quorum
	updates["version"] = bid.Version + 1
	return updates, nil
}
