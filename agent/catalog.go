package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	common_result "muidea.com/magicCommon/common"
	"muidea.com/magicCommon/model"
)

func (s *center) FetchCatalog(name string) (model.CatalogDetailView, bool) {
	type fetchResult struct {
		common_result.Result
		Catalog model.CatalogDetailView `json:"catalog"`
	}

	result := &fetchResult{}
	url := fmt.Sprintf("%s/%s?name=%s&authToken=%s&sessionID=%s", s.baseURL, "content/catalog/", name, s.authToken, s.sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.Catalog, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("fetch catalog failed, statusCode:%d", response.StatusCode)
		return result.Catalog, false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.Catalog, false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.Catalog, false
	}

	if result.ErrorCode == common_result.Success {
		return result.Catalog, true
	}

	log.Printf("fetch catalog failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.Catalog, false
}

func (s *center) QueryCatalog(catalogID int) (model.CatalogDetailView, bool) {
	type queryResult struct {
		common_result.Result
		Catalog model.CatalogDetailView `json:"catalog"`
	}

	result := &queryResult{}
	url := fmt.Sprintf("%s/%s/%d?authToken=%s&sessionID=%s", s.baseURL, "content/catalog", catalogID, s.authToken, s.sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.Catalog, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("query catalog failed, statusCode:%d", response.StatusCode)
		return result.Catalog, false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.Catalog, false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.Catalog, false
	}

	if result.ErrorCode == common_result.Success {
		return result.Catalog, true
	}

	log.Printf("query catalog failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.Catalog, false
}

func (s *center) CreateCatalog(name, description string, parent []model.Catalog, authToken, sessionID string) (model.SummaryView, bool) {
	type createParam struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Catalog     []model.Catalog `json:"catalog"`
	}

	type createResult struct {
		common_result.Result
		Catalog model.SummaryView `json:"catalog"`
	}

	param := createParam{Name: name, Description: description, Catalog: parent}
	result := &createResult{}

	data, err := json.Marshal(param)
	if err != nil {
		log.Printf("marshal create param failed, err:%s", err.Error())
		return result.Catalog, false
	}

	if len(authToken) == 0 {
		authToken = s.authToken
	}
	if len(sessionID) == 0 {
		sessionID = s.sessionID
	}

	bufferReader := bytes.NewBuffer(data)
	url := fmt.Sprintf("%s/%s?authToken=%s&sessionID=%s", s.baseURL, "content/catalog/", authToken, sessionID)
	request, err := http.NewRequest("POST", url, bufferReader)
	if err != nil {
		log.Printf("construct request failed, url:%s, err:%s", url, err.Error())
		return result.Catalog, false
	}

	request.Header.Set("content-type", "application/json")
	response, err := s.httpClient.Do(request)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.Catalog, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("create catalog failed, statusCode:%d", response.StatusCode)
		return result.Catalog, false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.Catalog, false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.Catalog, false
	}

	if result.ErrorCode == common_result.Success {
		return result.Catalog, true
	}

	log.Printf("create catalog failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.Catalog, false
}

func (s *center) DeleteCatalog(id int, authToken, sessionID string) bool {
	type deleteResult struct {
		common_result.Result
	}

	result := &deleteResult{}
	url := fmt.Sprintf("%s/%s/%d?authToken=%s&sessionID=%s", s.baseURL, "content/catalog", id, s.authToken, s.sessionID)
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Printf("construct request failed, url:%s, err:%s", url, err.Error())
		return false
	}

	request.Header.Set("content-type", "application/json")
	response, err := s.httpClient.Do(request)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("delete catalog failed, statusCode:%d", response.StatusCode)
		return false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return false
	}

	if result.ErrorCode == common_result.Success {
		return true
	}

	log.Printf("delete catalog failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return false
}
