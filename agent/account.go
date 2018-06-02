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

func (s *center) LoginAccount(account, password string) (model.AccountOnlineView, string, string, bool) {
	type loginParam struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}

	type loginResult struct {
		common_result.Result
		OnlineUser model.AccountOnlineView `json:"onlineUser"`
		SessionID  string                  `json:"sessionID"`
	}

	param := loginParam{Account: account, Password: password}
	result := &loginResult{}
	data, err := json.Marshal(param)
	if err != nil {
		log.Printf("marshal login param failed, err:%s", err.Error())
		return result.OnlineUser, "", "", false
	}

	bufferReader := bytes.NewBuffer(data)
	url := fmt.Sprintf("%s/%s", s.baseURL, "cas/user/")
	request, err := http.NewRequest("POST", url, bufferReader)
	if err != nil {
		log.Printf("construct request failed, url:%s, err:%s", url, err.Error())
		return result.OnlineUser, "", "", false
	}

	request.Header.Set("content-type", "application/json")
	response, err := s.httpClient.Do(request)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.OnlineUser, "", "", false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("login failed, statusCode:%d", response.StatusCode)
		return result.OnlineUser, "", "", false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.OnlineUser, "", "", false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.OnlineUser, "", "", false
	}

	if result.ErrorCode == common_result.Success {
		return result.OnlineUser, result.OnlineUser.AuthToken, result.SessionID, true
	}

	log.Printf("login failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.OnlineUser, "", "", false
}

func (s *center) LogoutAccount(authToken, sessionID string) bool {
	type logoutResult struct {
		common_result.Result
	}

	if len(authToken) == 0 || len(sessionID) == 0 {
		log.Print("illegal authToken or sessionID")
		return false
	}

	result := &logoutResult{}
	url := fmt.Sprintf("%s/%s/?authToken=%s&sessionID=%s", s.baseURL, "cas/user", authToken, sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("logout account failed, statusCode:%d", response.StatusCode)
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

	log.Printf("logout account failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return false
}

func (s *center) StatusAccount(authToken, sessionID string) (model.AccountOnlineView, bool) {
	type statusResult struct {
		common_result.Result
		OnlineUser model.AccountOnlineView `json:"onlineUser"`
		SessionID  string                  `json:"sessionID"`
	}

	result := &statusResult{}
	if len(authToken) == 0 || len(sessionID) == 0 {
		log.Print("illegal authToken or sessionID")
		return result.OnlineUser, false
	}

	url := fmt.Sprintf("%s/%s/?authToken=%s&sessionID=%s", s.baseURL, "cas/user", authToken, sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.OnlineUser, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("status account failed, statusCode:%d", response.StatusCode)
		return result.OnlineUser, false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.OnlineUser, false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.OnlineUser, false
	}

	if result.ErrorCode == common_result.Success {
		return result.OnlineUser, true
	}

	log.Printf("status account failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.OnlineUser, false
}
