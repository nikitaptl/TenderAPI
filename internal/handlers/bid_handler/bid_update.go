package bid_handler

import (
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators/bid_validators"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func UpdateBid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	updates, err := bid_validators.ValidateBidUpdates(r)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	vars := mux.Vars(r)
	bidId := vars["bidId"]
	bid, err := bid_validators.ValidateBidId(bidId)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	if bid, err = updateBid(bid, updates); err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	log.Printf("Предложение успешно обновлено: %v", bid)
	json.NewEncoder(w).Encode(bid)
}
