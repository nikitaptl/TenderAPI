package bid_handler

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"net/http"
)

func saveBidVersion(bid *models.Bid) error {
	bidVersion := models.BidVersion{
		BidID:       bid.ID,
		Version:     bid.Version,
		Name:        bid.Name,
		Description: bid.Description,
		Status:      bid.Status,
	}

	return database.DB.Create(&bidVersion).Error
}

// unsafe функция. Использовать только при выполненных валидациях
func updateBid(bid *models.Bid, updates map[string]interface{}) (*models.Bid, error) {
	if err := database.Update(bid, updates); err != nil {
		return nil, utils.NewErrorResponse("Неизвестная ошибка при обновлении тендера", http.StatusInternalServerError)
	}

	if err := saveBidVersion(bid); err != nil {
		return nil, utils.NewErrorResponse("Ошибка при сохранении версии тендера: "+err.Error(), http.StatusInternalServerError)
	}

	var bidResult models.Bid
	if err := database.DB.First(&bidResult, "id = ?", bid.ID).Error; err != nil {
		return nil, utils.NewErrorResponse("Ошибка при получении обновлённого тендера: "+err.Error(), http.StatusInternalServerError)
	}

	return &bidResult, nil
}
