package feedback_handler

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators/feedback_validator"
	"encoding/json"
	"net/http"
)

func CreateNewFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	feedback, err := feedback_validator.ValidateNewFeedback(r)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	if err = database.DB.Create(&feedback).Error; err != nil {
		utils.HandleHTTPError(w, utils.NewErrorResponse("Ошибка при создании отзыва", http.StatusInternalServerError))
		return
	}

	json.NewEncoder(w).Encode(feedback)
}
