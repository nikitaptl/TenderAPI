package bid_handler

import (
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators/bid_validators"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RollbackBid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	bidIDStr := vars["bidId"]
	versionStr := vars["version"]
	username := r.URL.Query().Get("username")

	bid, err := bid_validators.ValidateBidRollback(bidIDStr, versionStr, username)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}
	if err = saveBidVersion(bid); err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	log.Printf("Предложение успешно откачено до версии %s", versionStr)
	json.NewEncoder(w).Encode(bid)
}
