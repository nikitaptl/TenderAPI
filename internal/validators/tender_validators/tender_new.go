package tender_validators

import (
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/validators"
)

func ValidateNewTender(tender *models.Tender) error {
	if _, err := validators.OrganizationExist(tender.OrganizationID); err != nil {
		return err
	}
	employee, err := validators.UserExist(tender.CreatorUsername)
	if err != nil {
		return err
	}
	tender.Status = "Created"
	if err = validators.ValidateServiceType(tender.ServiceType); err != nil {
		return err
	}
	if _, err = validators.IsUserResponsible(tender.OrganizationID, employee.ID); err != nil {
		return err
	}
	return nil
}
