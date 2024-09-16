package tender_handler

import (
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators/tender_validators"
	"encoding/json"
	"log"
	"net/http"
)

func UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tender, err := tender_validators.ValidateStatusUpdate(r)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}
	updates := map[string]interface{}{"status": tender.Status,
		"version": tender.Version + 1}

	if tender, err = updateTender(tender, updates); err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	log.Printf("Статус тендера %s успешно изменён на %s", tender.ID, tender.Status)
	json.NewEncoder(w).Encode(tender)
}

func GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tender, err := tender_validators.ValidateGettingStatus(r)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	w.Write([]byte(`"` + tender.Status + `"`))
}
