package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	common_result "muidea.com/magicCommon/common"
	"muidea.com/magicCommon/model"
)

func (s *center) QueryMedia(id int) (model.MediaDetailView, bool) {
	type queryResult struct {
		common_result.Result
		Media model.MediaDetailView `json:"media"`
	}

	result := &queryResult{}
	url := fmt.Sprintf("%s/%s/%d?authToken=%s&sessionID=%s", s.baseURL, "content/media", id, s.authToken, s.sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.Media, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("query media failed, statusCode:%d", response.StatusCode)
		return result.Media, false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.Media, false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.Media, false
	}

	if result.ErrorCode == common_result.Success {
		return result.Media, true
	}

	log.Printf("query media failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.Media, false
}
