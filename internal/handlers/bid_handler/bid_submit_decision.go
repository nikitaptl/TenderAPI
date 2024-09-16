package bid_handler

import (
	"avitoTestTask/internal/database"
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"avitoTestTask/internal/validators"
	"avitoTestTask/internal/validators/bid_validators"
	"avitoTestTask/internal/validators/tender_validators"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func SubmitDecisionBid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidId := vars["bidId"]

	bid, err := bid_validators.ValidateBidId(bidId)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	request, err := bid_validators.ValidateSubmitDecision(r, bid)
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}

	var message string
	if request["decision"] == "reject" {
		err = rejectBid(bid)
		message = fmt.Sprintf("Предложение '%s' отклонено", bid.Name)
	} else {
		user, _ := validators.UserExist(request["username"])
		message, err = approveBid(bid, user)
	}
	if err != nil {
		utils.HandleHTTPError(w, err)
		return
	}
	log.Println(message)
	w.Write([]byte(message))
}

func rejectBid(bid *models.Bid) error {
	if err := database.DB.Model(&bid).Update("status", "CANCELED").Error; err != nil {
		return utils.NewErrorResponse("Ошибка при обновлении статуса предложения", http.StatusInternalServerError)
	}
	return nil
}

func approveBid(bid *models.Bid, user *models.Employee) (string, error) {
	if bid.RemainingApprovals > 1 {
		remainingApprovals := bid.RemainingApprovals - 1

		updates := map[string]interface{}{
			"approved_users":      append(bid.ApprovedUsers, user.ID),
			"remaining_approvals": remainingApprovals,
		}
		if err := database.DB.Model(&bid).Updates(updates).Error; err != nil {
			return "", utils.NewErrorResponse("Неизвестная ошибка при обновлении статуса предложения", http.StatusInternalServerError)
		}
		message := fmt.Sprintf("Предложение одобрено, осталось %d согласий", remainingApprovals)
		return message, nil
	}

	// Обновляем статус принятого предложения на "APPROVED"
	if err := database.DB.Model(&bid).Update("status", "Approved").Error; err != nil {
		return "", utils.NewErrorResponse("Ошибка при обновлении статуса предложения", http.StatusInternalServerError)
	}
	// Обновляем статус всех остальных предложений на "CANCELED"
	if err := database.DB.Model(&models.Bid{}).
		Where("tender_id = ? AND id != ?", bid.TenderID, bid.ID).
		Update("status", "Canceled").Error; err != nil {
		return "", utils.NewErrorResponse("Ошибка при обновлении статусов других предложений", http.StatusInternalServerError)
	}

	// Закрываем тендер
	tender, err := tender_validators.TenderExists(bid.TenderID)
	if err != nil {
		return "", err
	}
	if err = database.DB.Model(&tender).Update("status", "Closed").Error; err != nil {
		return "", utils.NewErrorResponse("Ошибка при обновлении статуса тендера", http.StatusInternalServerError)
	}

	message := "Предложение одобрено, тендер закрыт, остальные предложения отклонены"
	return message, nil
}
