package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"gotest.tools/v3/assert"
)

type crudTesting[M, R any] struct {
	modelName  string
	updater    Updater[M]
	comparer   Comparer[M, R]
	identifier Identifier[M]
}

func NewCrudTesting[M, R any](
	modelName string,
	updater Updater[M],
	comparer Comparer[M, R],
	identifier Identifier[M],
) *crudTesting[M, R] {
	return &crudTesting[M, R]{
		modelName:  modelName,
		updater:    updater,
		comparer:   comparer,
		identifier: identifier,
	}
}

type Updater[M any] interface {
	Update(M) M
}

type Comparer[M, R any] interface {
	Compare(M, R) bool
}

type Identifier[M any] interface {
	Identify(M) string
}

func (cTesting *crudTesting[M, R]) Execute(t *testing.T, httpUrl string, modelTest M) {

	client := &http.Client{}

	t.Run(fmt.Sprintf("should create %s successfully", cTesting.modelName), func(t *testing.T) {

		method := "POST"

		content, err := json.Marshal(modelTest)
		assert.NilError(t, err)

		payload := bytes.NewReader(content)

		req, err := http.NewRequest(method, httpUrl, payload)

		assert.NilError(t, err)

		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		assert.NilError(t, err)
		defer res.Body.Close()

		err = ReaderToStruct(res.Body, &modelTest)
		assert.NilError(t, err)

		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})

	t.Run(fmt.Sprintf("should read %s successfully", cTesting.modelName), func(t *testing.T) {

		method := "GET"

		req, err := http.NewRequest(method, httpUrl, nil)

		query := req.URL.Query()
		query.Add("id", cTesting.identifier.Identify(modelTest))
		req.URL.RawQuery = query.Encode()

		assert.NilError(t, err)

		res, err := client.Do(req)
		assert.NilError(t, err)
		defer res.Body.Close()

		var models []R
		err = ReaderToStruct(res.Body, &models)
		assert.NilError(t, err)
		assert.Equal(t, true, len(models) <= 100)

		var foundModel bool
		for _, model := range models {
			if cTesting.comparer.Compare(modelTest, model) {
				foundModel = true
			}
		}

		assert.Equal(t, true, foundModel)

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run(fmt.Sprintf("should update %s successfully", cTesting.modelName), func(t *testing.T) {

		method := "PUT"

		modelUpdate := cTesting.updater.Update(modelTest)

		content, err := json.Marshal(modelUpdate)
		assert.NilError(t, err)

		payload := bytes.NewReader(content)

		req, err := http.NewRequest(method, httpUrl, payload)

		assert.NilError(t, err)

		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		assert.NilError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run(fmt.Sprintf("should delete %s successfully", cTesting.modelName), func(t *testing.T) {

		method := "DELETE"

		parsedUrl, err := url.Parse(httpUrl)
		assert.NilError(t, err)

		query := parsedUrl.Query()
		query.Add("id", cTesting.identifier.Identify(modelTest))
		parsedUrl.RawQuery = query.Encode()

		req, err := http.NewRequest(method, parsedUrl.String(), nil)

		assert.NilError(t, err)

		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		assert.NilError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
