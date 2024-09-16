package bid_handler

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators/bid_validators"
	"encoding/json"
	"log"
	"net/http"
)

func NewBid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var bid models.Bid
	err := json.NewDecoder(r.Body).Decode(&bid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bid.Version = 1
	if err = bid_validators.ValidateNewBid(&bid); err != nil {
		utils.HandleHTTPError(w, err)
		return
	}
	if err = database.DB.Create(&bid).Error; err != nil {
		utils.HandleHTTPError(w, err)
	}
	if err = saveBidVersion(&bid); err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	log.Printf("Предложение успешно создано")
	json.NewEncoder(w).Encode(bid)
}
