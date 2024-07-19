package common

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

type CrudController[P, Q, T any] interface {
	Read(writer http.ResponseWriter, request *http.Request)
	Create(writer http.ResponseWriter, request *http.Request)
	Delete(writer http.ResponseWriter, request *http.Request)
	Update(writer http.ResponseWriter, request *http.Request)
}

type ReadFormatter interface {
	FormatAll(*http.Request, interface{}) (interface{}, error)
}

func NewCrudController[P, Q, T any](
	reader ReaderService[Q, T],
	inserter WriterService[P],
	updater WriterService[P],
	remover WriterService[T],
	queryCreator QueryCreator[Q],
	deleteCriteriaCreator QueryCreator[T],
	readFormatter ReadFormatter,
) CrudController[P, Q, T] {
	return &crudController[P, Q, T]{
		reader:                reader,
		inserter:              inserter,
		updater:               updater,
		remover:               remover,
		queryCreator:          queryCreator,
		deleteCriteriaCreator: deleteCriteriaCreator,
		readFormatter:         readFormatter,
	}
}

type QueryCreator[Q any] interface {
	Create(*http.Request) (*Q, error)
}

type crudController[P, Q, T any] struct {
	reader                ReaderService[Q, T]
	inserter              WriterService[P]
	updater               WriterService[P]
	remover               WriterService[T]
	queryCreator          QueryCreator[Q]
	deleteCriteriaCreator QueryCreator[T]
	readFormatter         ReadFormatter
}

func (controller *crudController[P, Q, T]) writeError(writer http.ResponseWriter, httpStatus int, err error) {
	log.Println(err)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpStatus)
	writer.Write(NewErrorToJson(err))
}

func (controller *crudController[P, Q, T]) Read(writer http.ResponseWriter, request *http.Request) {
	query, err := controller.queryCreator.Create(request)
	if err != nil {
		controller.writeError(writer, http.StatusInternalServerError, err)
		return
	}

	pagination := NewPagination()

	pageText := request.URL.Query().Get("page")

	if pageText != "" {
		page, err := strconv.ParseInt(pageText, 10, 32)
		if err != nil {
			controller.writeError(writer, http.StatusBadRequest, err)
			return
		}
		pagination.SetPage(int(page))
	}

	sizeText := request.URL.Query().Get("size")

	if sizeText != "" {
		size, err := strconv.ParseInt(sizeText, 10, 32)
		if err != nil {
			controller.writeError(writer, http.StatusBadRequest, err)
			return
		}
		pagination.SetSize(int(size))
	}

	list, err := controller.reader.Execute(request.Context(), query, pagination)

	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		controller.writeError(writer, http.StatusInternalServerError, err)
		return
	}

	if controller.readFormatter == nil {
		SendResponse(writer, list)
		return
	}

	formatted, err := controller.readFormatter.FormatAll(request, list)

	if err != nil {
		controller.writeError(writer, http.StatusInternalServerError, err)
		return
	}

	SendResponse(writer, formatted)
}

func (controller *crudController[P, Q, T]) Create(writer http.ResponseWriter, request *http.Request) {
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	payload := new(P)
	err = json.Unmarshal(bodyBytes, payload)
	if err != nil {
		controller.writeError(writer, http.StatusBadRequest, err)
		return
	}

	log.Println(payload)
	result, err := controller.inserter.Execute(request.Context(), payload)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		controller.writeError(writer, http.StatusInternalServerError, err)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	SendResponse(writer, result)
}

func (controller *crudController[P, Q, T]) Update(writer http.ResponseWriter, request *http.Request) {
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		controller.writeError(writer, http.StatusBadRequest, err)
		log.Println(err)
		return
	}
	payload := new(P)
	err = json.Unmarshal(bodyBytes, payload)
	if err != nil {
		controller.writeError(writer, http.StatusBadRequest, err)
		return
	}

	log.Println(payload)
	result, err := controller.updater.Execute(request.Context(), payload)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		controller.writeError(writer, http.StatusBadRequest, err)
		return
	}

	writer.WriteHeader(http.StatusOK)
	SendResponse(writer, result)
}

func (controller *crudController[P, Q, T]) Delete(writer http.ResponseWriter, request *http.Request) {
	query, err := controller.deleteCriteriaCreator.Create(request)

	if err != nil {
		controller.writeError(writer, http.StatusInternalServerError, err)
		return
	}

	_, err = controller.remover.Execute(request.Context(), query)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		controller.writeError(writer, http.StatusInternalServerError, err)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
