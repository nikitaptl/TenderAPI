package tender_validators

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var validStatuses = map[string]bool{
	"Created":   true,
	"Published": true,
	"Closed":    true,
}

var validStatusTransitions = map[string][]string{
	"Created":   {"Published", "Closed"}, // Created можно изменить на Published или Closed
	"Published": {"Closed"},              // Published можно изменить только на Closed
	"Closed":    {},                      // Closed нельзя изменить
}

func TenderExists(tenderId uuid.UUID) (*models.Tender, error) {
	var tender models.Tender
	result := database.DB.First(&tender, "id = ?", tenderId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		message := fmt.Sprintf("Тендер с ID %s не найден", tenderId)
		return nil, utils.NewErrorResponse(message, http.StatusNotFound)
	}
	if result.Error != nil {
		return nil, utils.NewErrorResponse("Неизвестная ошибка", http.StatusInternalServerError)
	}
	return &tender, nil
}

func ValidateTenderStatus(status string) error {
	if !validStatuses[status] {
		return utils.NewErrorResponse("Статус должен быть одним из: Created, Published, Closed", http.StatusBadRequest)
	}
	return nil
}

func ValidateNewTenderStatus(currentStatus, newStatus string) error {
	if err := ValidateTenderStatus(newStatus); err != nil {
		return err
	}
	validNextStatuses := validStatusTransitions[currentStatus]
	for _, status := range validNextStatuses {
		if status == newStatus {
			return nil
		}
	}
	return utils.NewErrorResponse("Нельзя изменить статус на "+newStatus+" из текущего "+currentStatus, http.StatusBadRequest)
}

func ValidateTenderVersion(tender *models.Tender, version string) (*models.TenderVersion, error) {
	versionInt, err := strconv.ParseUint(version, 10, 32)
	if err != nil {
		return nil, utils.NewErrorResponse("Неверный формат версии", http.StatusBadRequest)
	}

	var tenderVersion models.TenderVersion
	result := database.DB.First(&tenderVersion, "tender_id = ? AND version = ?", tender.ID, uint(versionInt))
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		message := fmt.Sprintf("Версия %d тендера %s не найдена", versionInt, tender.Name)
		return nil, utils.NewErrorResponse(message, http.StatusNotFound)
	}
	if result.Error != nil {
		return nil, utils.NewErrorResponse(result.Error.Error(), http.StatusInternalServerError)
	}

	return &tenderVersion, nil
}

func ValidateTenderId(tenderId string) (*models.Tender, error) {
	parsedTenderId, err := uuid.Parse(tenderId)
	if err != nil {
		return nil, utils.NewErrorResponse("Неверный формат ID тендера", http.StatusBadRequest)
	}
	tender, err := TenderExists(parsedTenderId)
	if err != nil {
		return nil, err
	}
	return tender, nil
}
