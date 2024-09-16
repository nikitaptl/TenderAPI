package feedback_validator

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators"
	"avitoTestTask/internal/validators/tender_validators"
	"github.com/gorilla/mux"
	"net/http"
)

func ValidateFeedbackList(r *http.Request) ([]models.Feedback, error) {
	vars := mux.Vars(r)
	tenderId := vars["tenderId"]
	creatorUsername := r.URL.Query().Get("authorUsername")
	requesterUsername := r.URL.Query().Get("requesterUsername")

	tender, err := tender_validators.ValidateTenderId(tenderId)
	if err != nil {
		return nil, err
	}
	_, err = validators.UserExist(creatorUsername)
	if err != nil {
		return nil, err
	}
	requester, err := validators.UserExist(requesterUsername)
	if err != nil {
		return nil, err
	}
	if _, err := validators.IsUserResponsible(tender.OrganizationID, requester.ID); err != nil {
		return nil, err
	}

	var bids []models.Bid
	if err = database.DB.Where("tender_id = ? AND creator_username = ?", tender.ID, creatorUsername).Find(&bids).Error; err != nil {
		return nil, utils.NewErrorResponse("Не удалось получить список предложений", http.StatusInternalServerError)
	}
	if len(bids) == 0 {
		return nil, utils.NewErrorResponse("Пользователь не участвовал в тендере", http.StatusBadRequest)
	}

	var feedbacks []models.Feedback
	if err := database.DB.Where("creator_username = ?", creatorUsername).Find(&feedbacks).Error; err != nil {
		return nil, utils.NewErrorResponse("Не удалось получить список отзывов", http.StatusInternalServerError)
	}
	return feedbacks, nil
}
