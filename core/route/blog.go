package route

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	cmsClient "github.com/muidea/magicCMS/client"
	cmsModel "github.com/muidea/magicCMS/model"
	commonCommon "github.com/muidea/magicCommon/common"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
)

func (s *Registry) filterPostList() interface{} {

	return nil
}

func (s *Registry) queryCatalog(catalog string, clnt cmsClient.Client) (ret []*cmsModel.CatalogLite, err error) {
	blogCatalog, blogErr := clnt.QueryCatalogTree(s.cmsCatalog, -1)
	if blogErr != nil {
		err = blogErr
		return
	}

	newCatalogItems := []string{}
	ret = []*cmsModel.CatalogLite{}

	catalogMapInfo := map[string]*cmsModel.CatalogLite{}
	for _, cv := range blogCatalog.Subs {
		catalogMapInfo[cv.Name] = cv.Lite()
	}

	items := strings.Split(catalog, ",")
	items = append(items, "current_catalog")
	for _, val := range items {
		cv, exist := catalogMapInfo[val]
		if !exist {
			catalogMapInfo[val] = nil
			newCatalogItems = append(newCatalogItems, val)
		} else if cv != nil {
			ret = append(ret, cv)
		}
	}

	for _, val := range newCatalogItems {
		newCatalog, newErr := clnt.CreateCatalog(val, "auto create blog catalog", blogCatalog.Lite())
		if newErr != nil {
			err = newErr
			return
		}

		ret = append(ret, newCatalog.Lite())
	}

	return
}

// PostBlog post blog
func (s *Registry) PostBlog(res http.ResponseWriter, req *http.Request) {
	type postParam struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Catalog string `json:"catalog"`
	}

	type postResult struct {
		commonDef.Result
		Redirect string `json:"redirect"`
	}

	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo, _ := curSession.GetOption(commonCommon.SessionIdentity)
	result := &postResult{}
	for {
		param := &postParam{}
		err := net.ParseJSONBody(req, param)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		if param.Title == "" || param.Catalog == "" {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数,输入参数为空"
			break
		}

		cmsClient := cmsClient.NewClient(s.cmsService)
		defer cmsClient.Release()

		cmsClient.BindSession(sessionInfo.(*commonCommon.SessionInfo))

		catalogList, catalogErr := s.queryCatalog(param.Catalog, cmsClient)
		if catalogErr != nil {
			log.Printf("queryCatalog failed, err:%s", catalogErr.Error())
			result.ErrorCode = commonDef.Failed
			result.Reason = "提交Blog失败, 查询分类出错"
			break
		}

		_, err = cmsClient.CreateArticle(param.Title, param.Content, catalogList)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "提交Blog失败, 保存出错"
			break
		}

		result.ErrorCode = commonDef.Success
		result.Redirect = "/"
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
