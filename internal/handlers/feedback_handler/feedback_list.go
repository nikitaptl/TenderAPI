package feedback_handler

import (
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators/feedback_validator"
	"encoding/json"
	"log"
	"net/http"
)

func GetFeedbackList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authorUsername := r.URL.Query().Get("authorUsername")

	feedbacks, err := feedback_validator.ValidateFeedbackList(r)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	log.Printf("Получен список из %d отзывов на автора %s", len(feedbacks), authorUsername)
	json.NewEncoder(w).Encode(feedbacks)
}
