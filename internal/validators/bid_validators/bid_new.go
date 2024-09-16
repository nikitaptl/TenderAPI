package bid_validators

import (
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators"
	"avitoTestTask/internal/validators/tender_validators"
	"net/http"
)

func ValidateNewBid(bid *models.Bid) error {
	tender, err := tender_validators.TenderExists(bid.TenderID)
	if err != nil {
		return err
	}
	if tender.Status != "Published" {
		return utils.NewErrorResponse("У вас нет доступа к данному тендеру, дождитесь публикации", http.StatusForbidden)
	}
	organization, err := validators.OrganizationExist(bid.OrganizationID)
	if err != nil {
		return err
	}
	employee, err := validators.UserExist(bid.CreatorUsername)
	if err != nil {
		return err
	}
	if _, err = validators.IsUserResponsible(organization.ID, employee.ID); err != nil {
		return err
	}

	bid.OrganizationTenderID = tender.OrganizationID
	bid.RemainingApprovals, err = FindQuorum(tender.OrganizationID)
	bid.Status = "Created"
	if err != nil {
		return err
	}

	return nil
}
