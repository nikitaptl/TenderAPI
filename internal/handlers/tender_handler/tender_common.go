package tender_handler

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"net/http"
)

func saveTenderVersion(tender *models.Tender) error {
	tenderVersion := models.TenderVersion{
		TenderID:    tender.ID,
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
		Status:      tender.Status,
		Version:     tender.Version,
	}
	return database.DB.Create(&tenderVersion).Error
}

// unsafe функция. Использовать только при выполненных валидациях
func updateTender(tender *models.Tender, updates map[string]interface{}) (*models.Tender, error) {
	if err := database.Update(tender, updates); err != nil {
		return nil, utils.NewErrorResponse("Неизвестная ошибка при обновлении тендера", http.StatusInternalServerError)
	}
	if err := saveTenderVersion(tender); err != nil {
		return nil, utils.NewErrorResponse("Неизвестная ошибка при сохранении версии тендера", http.StatusInternalServerError)
	}
	var tenderResult models.Tender
	if err := database.DB.First(&tenderResult, "id = ?", tender.ID).Error; err != nil {
		return nil, utils.NewErrorResponse("Неизвестная ошибка при получении обновлённого тендера", http.StatusInternalServerError)
	}
	return &tenderResult, nil
}
