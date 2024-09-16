package tender_handler

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators/tender_validators"
	"encoding/json"
	"log"
	"net/http"
)

func NewTender(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tender models.Tender
	err := json.NewDecoder(r.Body).Decode(&tender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = tender_validators.ValidateNewTender(&tender); err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	tender.Version = 1
	if err = database.DB.Create(&tender).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = saveTenderVersion(&tender); err != nil {
		http.Error(w, "Неизвестная ошибка при сохранении версии тендера", http.StatusInternalServerError)
		return
	}

	log.Printf("Тендер успешно создан: %v", tender)
	json.NewEncoder(w).Encode(tender)
}
