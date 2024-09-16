package bid_handler

import (
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators/bid_validators"
	"encoding/json"
	"log"
	"net/http"
)

func GetBidStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bid, err := bid_validators.ValidateGettingStatus(r)
	if err != nil {
		utils.HandleHTTPError(w, err)
	}
	w.Write([]byte(`"` + bid.Status + `"`))
}

func UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bid, err := bid_validators.ValidateStatusUpdate(r)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}
	updates := map[string]interface{}{"status": bid.Status,
		"version": bid.Version + 1}

	if bid, err = updateBid(bid, updates); err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	log.Printf("Статус предложения успешно обновлён: %v", bid)
	json.NewEncoder(w).Encode(bid)
}
