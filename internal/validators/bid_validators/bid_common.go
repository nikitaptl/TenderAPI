package bid_validators

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators"
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
	"Canceled":  true,
}

var validBidStatusTransitions = map[string][]string{
	"Created":   {"Published", "Canceled"},
	"Published": {"Canceled"},
	"Canceled":  {},
}

func BidExists(bidId uuid.UUID) (*models.Bid, error) {
	var bid models.Bid
	result := database.DB.First(&bid, "id = ?", bidId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		message := fmt.Sprintf("Предложение с ID %s не найдено", bidId)
		return nil, utils.NewErrorResponse(message, http.StatusNotFound)
	}

	if result.Error != nil {
		return nil, utils.NewErrorResponse("Неизвестная ошибка", http.StatusInternalServerError)
	}
	return &bid, nil
}

func ValidateBidId(bidId string) (*models.Bid, error) {
	parsedBidId, err := uuid.Parse(bidId)
	if err != nil {
		return nil, utils.NewErrorResponse("Неверный формат ID предложения", http.StatusBadRequest)
	}

	bid, err := BidExists(parsedBidId)
	if err != nil {
		return nil, err
	}

	return bid, nil
}

func ValidateBidVersion(bid *models.Bid, versionStr string) (*models.BidVersion, error) {
	version, err := strconv.ParseUint(versionStr, 10, 32)
	if err != nil {
		return nil, utils.NewErrorResponse("Неверный формат версии", http.StatusBadRequest)
	}

	var bidVersion models.BidVersion
	result := database.DB.Where("bid_id = ? AND version = ?", bid.ID, version).First(&bidVersion)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message := fmt.Sprintf("Версия %d для предложения с ID %s не найдена", version, bid.ID)
			return nil, utils.NewErrorResponse(message, http.StatusNotFound)
		}
		return nil, utils.NewErrorResponse("Ошибка при получении версии предложения", http.StatusInternalServerError)
	}

	return &bidVersion, nil
}

func ValidateNewBidStatus(currentStatus, newStatus string) error {
	if err := ValidateBidStatus(currentStatus); err != nil {
		return err
	}
	validNextStatuses, exists := validBidStatusTransitions[currentStatus]
	if !exists {
		return utils.NewErrorResponse("Недопустимый текущий статус: "+currentStatus, http.StatusBadRequest)
	}

	for _, status := range validNextStatuses {
		if status == newStatus {
			return nil
		}
	}
	return utils.NewErrorResponse("Нельзя изменить статус на "+newStatus+" из текущего "+currentStatus, http.StatusBadRequest)
}

func ValidateBidStatus(status string) error {
	if !validStatuses[status] {
		return utils.NewErrorResponse("Статус должен быть одним из: Created, Published, Closed", http.StatusBadRequest)
	}
	return nil
}

func FindQuorum(idTenderOrg uuid.UUID) (int, error) {
	var employeeCount int64
	if err := database.DB.Model(&models.OrganizationResponsible{}).
		Where("organization_id = ?", idTenderOrg).
		Count(&employeeCount).Error; err != nil {
		return 0, fmt.Errorf("Ошибка при подсчёте ответственных сотрудников: %v", err)
	}
	return int(min(3, employeeCount)), nil
}

func ValidateBidAndUser(bidIdStr string, username string) (*models.Bid, error) {
	user, err := validators.UserExist(username)
	if err != nil {
		return nil, err
	}
	bid, err := ValidateBidId(bidIdStr)
	if err != nil {
		return nil, err
	}
	if _, err = validators.IsUserResponsible(bid.OrganizationID, user.ID); err != nil {
		return nil, err
	}
	return bid, nil
}
