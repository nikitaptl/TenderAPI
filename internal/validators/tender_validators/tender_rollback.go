package tender_validators

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators"
	"net/http"
)

func ValidateTenderRollback(tenderIDStr string, versionStr string, username string) (*models.Tender, error) {
	user, err := validators.UserExist(username)
	if err != nil {
		return nil, err
	}
	tender, err := ValidateTenderId(tenderIDStr)
	if err != nil {
		return nil, err
	}
	if _, err := validators.IsUserResponsible(tender.OrganizationID, user.ID); err != nil {
		return nil, err
	}
	tenderVersion, err := ValidateTenderVersion(tender, versionStr)
	if err != nil {
		return nil, err
	}
	result := database.DB.Model(&models.Tender{}).Where("id = ?", tender.ID).Updates(models.Tender{
		Name:        tenderVersion.Name,
		Description: tenderVersion.Description,
		ServiceType: tenderVersion.ServiceType,
		Status:      tenderVersion.Status,
		Version:     tender.Version + 1,
	})
	if result.Error != nil {
		return nil, utils.NewErrorResponse("Ошибка при обновлении тендера", http.StatusInternalServerError)
	}
	if err = database.DB.First(tender, "id = ?", tender.ID).Error; err != nil {
		return nil, utils.NewErrorResponse("Ошибка при получении информации об обновлённом тендере", http.StatusInternalServerError)
	}
	return tender, nil
}
