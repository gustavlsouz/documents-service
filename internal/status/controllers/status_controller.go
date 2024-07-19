package controllers

import (
	"net/http"

	"github.com/gustavlsouz/documents-service/internal/common"
	"github.com/gustavlsouz/documents-service/internal/status/services"
)

type statusController struct {
	GetStatusService services.GetStatusService
}

func NewStatusController() *statusController {
	return &statusController{
		GetStatusService: services.NewGetStatusService(),
	}
}

func (controller *statusController) GetStatus(writer http.ResponseWriter, request *http.Request) {
	status := controller.GetStatusService.Execute(request.Context())
	common.SendResponse(writer, status)
}
