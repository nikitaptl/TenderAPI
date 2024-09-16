package bid_validators

import (
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators"
	"net/http"
)

func ValidateSubmitDecision(r *http.Request, bid *models.Bid) (map[string]string, error) {
	username := r.URL.Query().Get("username")
	user, err := validators.UserExist(username)
	if err != nil {
		return nil, err
	}

	decision := r.URL.Query().Get("decision")
	if decision != "approve" && decision != "reject" {
		return nil, utils.NewErrorResponse("Поле 'decision' должно быть 'approve' или 'reject'", http.StatusBadRequest)
	}

	if _, err = validators.IsUserResponsible(bid.OrganizationTenderID, user.ID); err != nil {
		return nil, err
	}
	for _, userId := range bid.ApprovedUsers {
		if userId == user.ID && decision == "approve" {
			return nil, utils.NewErrorResponse("Пользователь уже одобрил это предложение", http.StatusConflict)
		}
	}
	request := map[string]string{
		"username": username,
		"decision": decision,
	}
	return request, nil
}
