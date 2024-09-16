package validators

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
	"strings"
)

var validServiceTypes = map[string]bool{
	"Construction": true,
	"Delivery":     true,
	"Manufacture":  true,
}

func OrganizationExist(organizationId uuid.UUID) (*models.Organization, error) {
	var organization models.Organization
	result := database.DB.First(&organization, "id = ?", organizationId)
	if result.Error == gorm.ErrRecordNotFound {
		message := fmt.Sprintf("Организация с id = %s не найдена", organizationId)
		return nil, utils.NewErrorResponse(message, http.StatusNotFound)
	}
	if result.Error != nil {
		return nil, utils.NewErrorResponse("Ошибка при поиске организации", http.StatusInternalServerError)
	}
	return &organization, nil
}

func UserExist(employeeUsername string) (*models.Employee, error) {
	if employeeUsername == "" {
		return nil, utils.NewErrorResponse("Вы не ввели username", http.StatusUnauthorized)
	}

	var employee models.Employee
	result := database.DB.First(&employee, "username = ?", employeeUsername)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		message := fmt.Sprintf("Пользователь '%s' не существует или некорректен", employeeUsername)
		return nil, utils.NewErrorResponse(message, http.StatusUnauthorized)
	}
	if result.Error != nil {
		return nil, utils.NewErrorResponse("Ошибка при поиске пользователя", http.StatusInternalServerError)
	}
	return &employee, nil
}

func IsUserResponsible(organizationId uuid.UUID, employeeId uuid.UUID) (*models.OrganizationResponsible, error) {
	var orgResp models.OrganizationResponsible
	err := database.DB.Where("organization_id = ? AND user_id = ?", organizationId, employeeId).First(&orgResp).Error
	if err == gorm.ErrRecordNotFound {
		message := fmt.Sprintf("Пользователь %s не ответственнен за организацию %s", employeeId, organizationId)
		return nil, utils.NewErrorResponse(message, http.StatusForbidden)
	}
	if err != nil {
		return nil, utils.NewErrorResponse("Ошибка при проверке ответственности пользователя за организацию", http.StatusInternalServerError)
	}
	return &orgResp, nil
}

func ValidateLimit(r *http.Request) (int, error) {
	limit := 5
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			return 0, utils.NewErrorResponse("Неверный параметр 'limit'", http.StatusBadRequest)
		}
		limit = parsedLimit
	}
	return limit, nil
}

func ValidateOffset(r *http.Request) (int, error) {
	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			return 0, utils.NewErrorResponse("Неверный параметр 'offset'", http.StatusBadRequest)
		}
		offset = parsedOffset
	}
	return offset, nil
}

func ValidateServiceTypes(r *http.Request) ([]string, error) {
	serviceTypesParam := r.URL.Query().Get("service_type")

	if serviceTypesParam == "" {
		return nil, nil
	}

	serviceTypes := strings.Split(serviceTypesParam, ",")
	for _, serviceType := range serviceTypes {
		if !validServiceTypes[serviceType] {
			return nil, utils.NewErrorResponse("Неверный параметр 'service_type'", http.StatusBadRequest)
		}
	}
	return serviceTypes, nil
}

func ValidateServiceType(serviceType string) error {
	if !validServiceTypes[serviceType] {
		return utils.NewErrorResponse("Неверный параметр 'service_type'", http.StatusBadRequest)
	}
	return nil
}
