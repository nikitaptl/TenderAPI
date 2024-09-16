package tender_handler

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators"
	"encoding/json"
	"log"
	"net/http"
)

func Tenders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Параметры для запросов с пагинацией
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

	serviceTypes, err := validators.ValidateServiceTypes(r)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	query := database.DB.Model(&models.Tender{}).Limit(limit).Offset(offset).Order("name ASC")
	if len(serviceTypes) > 0 {
		query = query.Where("service_type IN (?)", serviceTypes)
	}

	var tenders []models.Tender
	if err = query.Find(&tenders).Error; err != nil {
		utils.HandleHTTPError(w, utils.NewErrorResponse("Ошибка при получении списка тендеров", http.StatusInternalServerError))
		return
	}

	log.Printf("Передан список из %d тендеров", len(tenders))
	json.NewEncoder(w).Encode(tenders)
}

func TendersMy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := r.URL.Query().Get("username")
	if username == "" {
		utils.HandleHTTPError(w, utils.NewErrorResponse("Не передан параметр 'username'", http.StatusBadRequest))
		return
	}
	if _, err := validators.UserExist(username); err != nil {
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

	var tenders []models.Tender
	query := database.DB.Where("creator_username = ?", username)

	err = query.Limit(limit).Offset(offset).Order("name ASC").Find(&tenders).Error
	if err != nil {
		utils.HandleHTTPError(w, utils.NewErrorResponse("Ошибка при получении тендеров", http.StatusInternalServerError))
		return
	}

	log.Printf("Пользователю %s передан список из %d тендеров", username, len(tenders))
	json.NewEncoder(w).Encode(tenders)
}
