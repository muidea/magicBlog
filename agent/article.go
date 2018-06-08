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

func (s *center) QueryArticle(id int) (model.ArticleDetailView, bool) {
	type queryResult struct {
		common_result.Result
		Article model.ArticleDetailView `json:"article"`
	}

	result := &queryResult{}
	url := fmt.Sprintf("%s/%s/%d?authToken=%s&sessionID=%s", s.baseURL, "content/article", id, s.authToken, s.sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.Article, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("query article failed, statusCode:%d", response.StatusCode)
		return result.Article, false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.Article, false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.Article, false
	}

	if result.ErrorCode == common_result.Success {
		return result.Article, true
	}

	log.Printf("query article failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.Article, false
}

func (s *center) CreateArticle(title, content string, catalog []model.Catalog, authToken, sessionID string) (model.SummaryView, bool) {
	type createParam struct {
		Name    string          `json:"name"`
		Content string          `json:"content"`
		Catalog []model.Catalog `json:"catalog"`
	}

	type createResult struct {
		common_result.Result
		Article model.SummaryView `json:"article"`
	}

	param := createParam{Name: title, Content: content, Catalog: catalog}
	result := &createResult{}
	data, err := json.Marshal(param)
	if err != nil {
		log.Printf("marshal create param failed, err:%s", err.Error())
		return result.Article, false
	}

	if len(authToken) == 0 {
		authToken = s.authToken
	}
	if len(sessionID) == 0 {
		sessionID = s.sessionID
	}

	bufferReader := bytes.NewBuffer(data)
	url := fmt.Sprintf("%s/%s?authToken=%s&sessionID=%s", s.baseURL, "content/article/", authToken, sessionID)
	request, err := http.NewRequest("POST", url, bufferReader)
	if err != nil {
		log.Printf("construct request failed, url:%s, err:%s", url, err.Error())
		return result.Article, false
	}

	request.Header.Set("content-type", "application/json")
	response, err := s.httpClient.Do(request)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.Article, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("create article failed, statusCode:%d", response.StatusCode)
		return result.Article, false
	}

	contentData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.Article, false
	}

	err = json.Unmarshal(contentData, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.Article, false
	}

	if result.ErrorCode == common_result.Success {
		return result.Article, true
	}

	log.Printf("create article failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.Article, false
}

func (s *center) DeleteArticle(id int, authToken, sessionID string) bool {
	type deleteResult struct {
		common_result.Result
	}

	result := &deleteResult{}
	url := fmt.Sprintf("%s/%s/%d?authToken=%s&sessionID=%s", s.baseURL, "content/article", id, s.authToken, s.sessionID)
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
		log.Printf("delete article failed, statusCode:%d", response.StatusCode)
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

	log.Printf("query article failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return false
}
