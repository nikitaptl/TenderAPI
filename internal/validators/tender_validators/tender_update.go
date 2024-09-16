package tender_validators

import (
	"avitoTestTask/internal/models"
	"avitoTestTask/internal/utils"
	"encoding/json"
	"io"
	"net/http"
)

func ValidateUpdates(body io.ReadCloser) (map[string]interface{}, error) {
	defer body.Close()

	var request models.UpdateTenderRequest
	if err := json.NewDecoder(body).Decode(&request); err != nil {
		return nil, utils.NewErrorResponse(err.Error(), http.StatusBadRequest)
	}

	updates := make(map[string]interface{})
	if request.Name != nil {
		updates["name"] = *request.Name
	}
	if request.Description != nil {
		updates["description"] = *request.Description
	}
	if request.ServiceType != nil {
		updates["service_type"] = *request.ServiceType
	}
	if len(updates) == 0 {
		return nil, utils.NewErrorResponse("Не переданы поля, которые могут быть обновлены", http.StatusBadRequest)
	}
	return updates, nil
}
