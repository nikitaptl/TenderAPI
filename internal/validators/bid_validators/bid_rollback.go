package bid_validators

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"net/http"
)

func ValidateBidRollback(bidIDStr string, versionStr string, username string) (*models.Bid, error) {
	bid, err := ValidateBidAndUser(bidIDStr, username)
	if err != nil {
		return nil, err
	}

	bidVersion, err := ValidateBidVersion(bid, versionStr)
	if err != nil {
		return nil, err
	}

	result := database.DB.Model(&models.Bid{}).Where("id = ?", bid.ID).Updates(models.Bid{
		Name:        bidVersion.Name,
		Description: bidVersion.Description,
		Status:      bidVersion.Status,
		Version:     bid.Version + 1,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if err = database.DB.First(bid, "id = ?", bid.ID).Error; err != nil {
		return nil, utils.NewErrorResponse("Ошибка при получении информации об обновлённом предложении", http.StatusInternalServerError)
	}
	return bid, nil
}
