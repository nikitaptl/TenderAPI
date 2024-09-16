package bid_handler

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators"
	"avitoTestTask/internal/validators/tender_validators"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func BidsMy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := r.URL.Query().Get("username")
	if username == "" {
		utils.HandleHTTPError(w, utils.NewErrorResponse("Не передан параметр 'username'", http.StatusBadRequest))
		return
	}
	if _, err := validators.UserExist(username); err != nil {
		utils.HandleHTTPError(w, utils.NewErrorResponse("Пользователь не существует или некорректен", http.StatusUnauthorized))
		return
	}

	limit, err := validators.ValidateLimit(r)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}
	offset, err := validators.ValidateOffset(r)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	var bids []models.Bid
	query := database.DB.Where("creator_username = ?", username).
		Order("name ASC").
		Limit(limit).
		Offset(offset)

	if err = query.Find(&bids).Error; err != nil {
		utils.HandleHTTPError(w, utils.NewErrorResponse("Ошибка при получении предложений", http.StatusInternalServerError))
		return
	}

	log.Printf("Пользователю %s передан список из %d его предложений", username, len(bids))
	json.NewEncoder(w).Encode(bids)
}

func BidsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	tenderId := vars["tenderId"]
	username := r.URL.Query().Get("username")

	user, err := validators.UserExist(username)
	if err != nil {
		utils.HandleHTTPError(w, utils.NewErrorResponse("Пользователь не существует или некорректен", http.StatusUnauthorized))
		return
	}
	tender, err := tender_validators.ValidateTenderId(tenderId)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}
	if _, err = validators.IsUserResponsible(tender.OrganizationID, user.ID); err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	limit, err := validators.ValidateLimit(r)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}
	offset, err := validators.ValidateOffset(r)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	var bids []models.Bid
	if err = database.DB.Where("tender_id = ? AND status = ?", tender.ID, "PUBLISHED").
		Limit(limit).Offset(offset).
		Order("name ASC").
		Find(&bids).Error; err != nil {
		utils.HandleHTTPError(w, utils.NewErrorResponse("Ошибка при получении предложений", http.StatusInternalServerError))
		return
	}

	log.Printf("Тендеру '%s' соответствует %d опубликованных предложений, список передан", tender.Name, len(bids))
	json.NewEncoder(w).Encode(bids)
}
