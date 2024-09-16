package feedback_validator

import (
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators"
	"avitoTestTask/internal/validators/bid_validators"
	"github.com/gorilla/mux"
	"net/http"
)

func ValidateNewFeedback(r *http.Request) (*models.Feedback, error) {
	vars := mux.Vars(r)
	bidIdStr := vars["bidId"]

	bidFeedback := r.URL.Query().Get("bidFeedback")
	username := r.URL.Query().Get("username")

	bid, err := bid_validators.ValidateBidId(bidIdStr)
	if err != nil {
		return nil, err
	}
	if bid.Status != "APPROVED" {
		return nil, utils.NewErrorResponse("Нельзя оставить отзыв на неодобренное предложение", http.StatusBadRequest)
	}

	user, err := validators.UserExist(username)
	if err != nil {
		return nil, err
	}
	if _, err = validators.IsUserResponsible(bid.OrganizationTenderID, user.ID); err != nil {
		return nil, err
	}
	return &models.Feedback{
		BidID:           bid.ID,
		CreatorUsername: bid.CreatorUsername,
		FeedbackText:    bidFeedback,
	}, nil
}
