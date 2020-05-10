package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/muidea/magicBlog/config"
	"github.com/muidea/magicBlog/model"
	cmsClient "github.com/muidea/magicCMS/client"
	cmsModel "github.com/muidea/magicCMS/model"
	casModel "github.com/muidea/magicCas/model"
	commonCommon "github.com/muidea/magicCommon/common"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/foundation/util"
)

const currentCatalog = "current_catalog"
const archiveCatalog = "archive_catalog"
const systemCatalog = "system_catalog"
const authorCatalog = "author_catalog"
const settingTitle = "setting"

type archiveBlogTask struct {
	registry     *Registry
	preTimeStamp *time.Time
}

func (s *archiveBlogTask) Run() {
	current := time.Now()
	if s.preTimeStamp == nil {
		s.preTimeStamp = &current
		if current.Day() > 1 {
			// 不是每月的第一天，不用计算
			return
		}
	} else {
		if s.preTimeStamp.Month() == current.Month() {
			s.preTimeStamp = &current
			// 月份没有变化，也不用计算
			return
		}
	}

	preTime := *s.preTimeStamp
	s.preTimeStamp = &current

	log.Printf("archive blog....., date:%s", preTime.Format("2006-01-02"))
	s.registry.archiveBlog()
}

func (s *Registry) confirmEndpoint() (ret *commonCommon.SessionInfo, err error) {
	cmsClient := cmsClient.NewClient(s.cmsService)
	defer cmsClient.Release()

	identityID := config.IdentityID()
	authToken := config.AuthToken()
	_, ret, err = cmsClient.ConfirmEndpoint(identityID, authToken)

	return
}

func (s *Registry) getCMSClient() (ret cmsClient.Client, err error) {
	sessionInfo, sessionErr := s.confirmEndpoint()
	if sessionErr != nil {
		log.Printf("confirmEndpoint failed, err:%s", sessionErr.Error())
		err = sessionErr
		return
	}

	sessionInfo.Scope = commonCommon.ShareSession
	cmsClient := cmsClient.NewClient(s.cmsService)
	cmsClient.BindSession(sessionInfo)

	ret = cmsClient

	return
}

func (s *Registry) queryBlogCommon(clnt cmsClient.Client) (catalogs []*cmsModel.CatalogLite, archives []*cmsModel.CatalogLite, articleList []*cmsModel.ArticleView, err error) {
	blogCatalog, blogErr := clnt.QueryCatalogTree(s.cmsCatalog, 2)
	if blogErr != nil {
		err = blogErr
		return
	}

	var archiveTree *cmsModel.CatalogTree
	var systemTree *cmsModel.CatalogTree
	var authorTree *cmsModel.CatalogTree
	catalogs = []*cmsModel.CatalogLite{}
	archives = []*cmsModel.CatalogLite{}
	for _, cv := range blogCatalog.Subs {
		switch cv.Name {
		case currentCatalog:
			s.currentCatalog = cv.Lite()
		case archiveCatalog:
			archiveTree = cv
		case systemCatalog:
			systemTree = cv
		case authorCatalog:
			authorTree = cv
		default:
			catalogs = append(catalogs, cv.Lite())
		}
	}

	if s.currentCatalog == nil {
		catalogPtr, catalogErr := clnt.CreateCatalog(currentCatalog, "auto create current catalog", blogCatalog.Lite())
		if catalogErr != nil {
			err = catalogErr
			return
		}

		s.currentCatalog = catalogPtr.Lite()
	}

	if archiveTree != nil {
		s.archiveCatalog = archiveTree.Lite()

		for _, cv := range archiveTree.Subs {
			archives = append(archives, cv.Lite())
		}
	} else {
		catalogPtr, catalogErr := clnt.CreateCatalog(archiveCatalog, "auto create archive catalog", blogCatalog.Lite())
		if catalogErr != nil {
			err = catalogErr
			return
		}

		s.archiveCatalog = catalogPtr.Lite()
	}

	if systemTree != nil {
		s.systemCatalog = systemTree.Lite()

		articleList, err = s.queryArticleList(clnt, systemTree.Lite(), nil)
	} else {
		catalogPtr, catalogErr := clnt.CreateCatalog(systemCatalog, "auto create system catalog", blogCatalog.Lite())
		if catalogErr != nil {
			err = catalogErr
			return
		}

		s.systemCatalog = catalogPtr.Lite()
	}

	if authorTree == nil {
		catalogPtr, catalogErr := clnt.CreateCatalog(authorCatalog, "auto create author catalog", blogCatalog.Lite())
		if catalogErr != nil {
			err = catalogErr
			return
		}

		s.authorCatalog = catalogPtr.Lite()
	} else {
		s.authorCatalog = authorTree.Lite()
	}

	return
}

func (s *Registry) getBlogSetting(articleList []*cmsModel.ArticleView) (ret *model.Setting, err error) {
	var settingPtr *cmsModel.ArticleView
	for _, val := range articleList {
		if val.Title == settingTitle {
			settingPtr = val
			break
		}
	}

	ret = &model.Setting{}
	if settingPtr != nil {
		err = json.Unmarshal([]byte(settingPtr.Content), ret)
		if err == nil {
			ret.ID = settingPtr.ID
		}
	}

	return
}

func (s *Registry) queryBlogPost(filter *filter, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	articleList, articleErr := s.queryArticleList(clnt, s.currentCatalog, filter.pageFilter)
	if articleErr != nil {
		err = articleErr
		return
	}

	var articlePtr *cmsModel.ArticleView
	for _, val := range articleList {
		fileName := fmt.Sprintf("%s.html", val.Title)
		if val.ID == filter.postID && fileName == filter.fileName {
			articlePtr = val
			break
		}
	}

	if articlePtr == nil {
		fileName = "404.html"
		return
	}

	info := map[string]interface{}{}
	commentList, commentErr := s.queryComments(clnt, articlePtr.ID, filter.pageFilter)
	if commentErr == nil {
		info["Comments"] = commentList
	}
	info["Content"] = articlePtr

	fileName = "post.html"
	content = info
	return
}

func (s *Registry) filterBlogArchive(filter *filter, archives []*cmsModel.CatalogLite, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	var archivePtr *cmsModel.CatalogLite
	for _, val := range archives {
		if filter.archiveName == val.Name {
			archivePtr = val
			break
		}
	}
	if archivePtr == nil {
		err = fmt.Errorf("illegal archive, name:%s", filter.archiveName)
		return
	}
	if filter.fileName != "" {
		articlePtr, articleErr := s.queryArticle(clnt, filter.postID)
		if articleErr != nil {
			err = articleErr
			return
		}

		info := map[string]interface{}{}
		commentList, commentErr := s.queryComments(clnt, articlePtr.ID, filter.pageFilter)
		if commentErr == nil {
			info["Comments"] = commentList
		}
		info["Content"] = articlePtr

		fileName = "post.html"
		content = info
		return
	}

	articleList, articleErr := s.queryArticleList(clnt, archivePtr, filter.pageFilter)
	if articleErr != nil {
		err = articleErr
		return
	}

	fileName = "index.html"
	content = articleList
	return
}

func (s *Registry) filterBlogCatalog(filter *filter, catalogs []*cmsModel.CatalogLite, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	var catalogPtr *cmsModel.CatalogLite
	for _, val := range catalogs {
		if filter.catalogName == val.Name {
			catalogPtr = val
			break
		}
	}
	if catalogPtr == nil {
		err = fmt.Errorf("illegal catalog, name:%s", filter.catalogName)
		return
	}

	if filter.fileName != "" {
		articlePtr, articleErr := s.queryArticle(clnt, filter.postID)
		if articleErr != nil {
			err = articleErr
			return
		}

		info := map[string]interface{}{}
		commentList, commentErr := s.queryComments(clnt, articlePtr.ID, filter.pageFilter)
		if commentErr == nil {
			info["Comments"] = commentList
		}
		info["Content"] = articlePtr

		fileName = "post.html"
		content = info
		return
	}

	articleList, articleErr := s.queryArticleList(clnt, catalogPtr, filter.pageFilter)
	if articleErr != nil {
		err = articleErr
		return
	}

	fileName = "index.html"
	content = articleList
	return
}

func (s *Registry) filterBlogAuthor(filter *filter, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	if s.authorCatalog == nil {
		err = fmt.Errorf("illegal author catalog")
		return
	}

	articleList, articleErr := s.queryArticleList(clnt, s.authorCatalog, nil)
	if articleErr != nil {
		err = articleErr
		return
	}

	var articlePtr *cmsModel.ArticleView
	for _, val := range articleList {
		fileName := fmt.Sprintf("%s.html", val.Title)
		if fileName == filter.fileName {
			articlePtr = val
			break
		}
	}
	if articlePtr == nil {
		fileName = "404.html"
		return
	}

	info := map[string]interface{}{}
	if articlePtr != nil {
		commentList, commentErr := s.queryComments(clnt, articlePtr.ID, filter.pageFilter)
		if commentErr == nil {
			info["Comments"] = commentList
		}
		info["Content"] = articlePtr
	}

	fileName = "post.html"
	content = info
	return
}

func (s *Registry) queryBlogLogin(filter *filter, articles []*cmsModel.ArticleView, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	var articlePtr *cmsModel.ArticleView
	for _, val := range articles {
		fileName := fmt.Sprintf("%s.html", val.Title)
		if fileName == filter.fileName {
			articlePtr = val
			break
		}
	}

	fileName = "login.html"
	content = articlePtr

	return
}

func (s *Registry) queryBlogAbout(filter *filter, articles []*cmsModel.ArticleView, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	var articlePtr *cmsModel.ArticleView
	for _, val := range articles {
		fileName := fmt.Sprintf("%s.html", val.Title)
		if fileName == filter.fileName {
			articlePtr = val
			break
		}
	}

	if articlePtr == nil {
		fileName = "404.html"
		return
	}

	info := map[string]interface{}{}
	info["Content"] = articlePtr

	fileName = "about.html"
	content = info

	return
}

func (s *Registry) queryBlogContact(filter *filter, articles []*cmsModel.ArticleView, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	var articlePtr *cmsModel.ArticleView
	for _, val := range articles {
		fileName := fmt.Sprintf("%s.html", val.Title)
		if fileName == filter.fileName {
			articlePtr = val
			break
		}
	}

	if articlePtr == nil {
		fileName = "404.html"
		return
	}

	info := map[string]interface{}{}
	commentList, commentErr := s.queryComments(clnt, articlePtr.ID, filter.pageFilter)
	if commentErr == nil {
		info["Comments"] = commentList
	}
	info["Content"] = articlePtr

	fileName = "contact.html"
	content = info
	return
}

func (s *Registry) filterBlogPostList(filter *filter, clnt cmsClient.Client) (ret []*cmsModel.ArticleView, err error) {
	articleList, articleErr := s.queryArticleList(clnt, s.currentCatalog, filter.pageFilter)
	if articleErr != nil {
		err = articleErr
		return
	}

	ret = articleList
	return
}

func (s *Registry) queryBlogPostEdit(filter *filter, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	if filter.action != "update_post" || filter.postID <= 0 {
		err = fmt.Errorf("illegal action, action:%s, postID:%d", filter.action, filter.postID)
		return
	}

	articleView, articleErr := s.queryArticle(clnt, filter.postID)
	if articleErr != nil {
		err = articleErr
		return
	}

	info := map[string]interface{}{}
	info["ID"] = articleView.ID
	info["Title"] = articleView.Title
	info["Content"] = articleView.Content

	catalogs := ""
	for _, val := range articleView.Catalog {
		if val.Name == s.currentCatalog.Name {
			continue
		}

		catalogs = fmt.Sprintf("%s,%s", catalogs, val.Name)
	}
	info["Catalog"] = strings.Trim(catalogs, ",")

	content = info
	fileName = "edit.html"
	return
}

func (s *Registry) queryBlogSetting(filter *filter, articles []*cmsModel.ArticleView, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	fileName = "setting.html"
	return
}

func (s *Registry) deleteBlogPost(filter *filter, clnt cmsClient.Client) (err error) {
	if filter.action != "delete_post" || filter.postID <= 0 {
		err = fmt.Errorf("illegal action, action:%s, postID:%d", filter.action, filter.postID)
		return
	}

	_, err = s.deleteArticle(clnt, filter.postID)
	return
}

func (s *Registry) deleteBlogCatalog(filter *filter, catalogs []*cmsModel.CatalogLite, clnt cmsClient.Client) (err error) {
	if filter.action != "delete_catalog" || filter.catalogID <= 0 {
		err = fmt.Errorf("illegal action, action:%s, catalogId:%d", filter.action, filter.catalogID)
		return
	}

	var catalogPtr *cmsModel.CatalogLite
	for _, val := range catalogs {
		if val.ID == filter.catalogID {
			catalogPtr = val
			break
		}
	}
	if catalogPtr == nil {
		err = fmt.Errorf("illegal action, action:%s, catalogId:%d", filter.action, filter.catalogID)
		return
	}

	articleList, articleErr := s.queryArticleList(clnt, catalogPtr, nil)
	if articleErr != nil {
		err = articleErr
		return
	}
	if len(articleList) > 0 {
		err = fmt.Errorf("delete catalog failed, catalog is busy. action:%s, catalogId:%d", filter.action, filter.catalogID)
		return
	}

	_, err = s.deleteCatalog(clnt, filter.catalogID)
	return
}

func (s *Registry) queryArticle(clnt cmsClient.Client, id int) (ret *cmsModel.ArticleView, err error) {
	blogArticle, blogErr := clnt.QueryArticle(id)
	if blogErr != nil {
		err = blogErr
		log.Printf("QueryArticle failed, err:%s", err.Error())
		return
	}

	ret = blogArticle
	return
}

func (s *Registry) queryArticleList(clnt cmsClient.Client, catalog *cmsModel.CatalogLite, pageFilter *util.PageFilter) (ret []*cmsModel.ArticleView, err error) {
	if pageFilter == nil {
		pageFilter = &util.PageFilter{PageSize: 10, PageNum: 1}
	}

	blogArticle, _, blogErr := clnt.FilterArticle(catalog, pageFilter)
	if blogErr != nil {
		err = blogErr
		log.Printf("FilterArticle failed, err:%s", err.Error())
		return
	}

	ret = blogArticle
	return
}

func (s *Registry) getCatalogs(catalog string, clnt cmsClient.Client) (ret []*cmsModel.CatalogLite, err error) {
	blogCatalog, blogErr := clnt.QueryCatalogTree(s.cmsCatalog, 2)
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

	items := strings.Split(strings.Trim(catalog, " "), ",")
	for _, val := range items {
		if val == currentCatalog || val == archiveCatalog {
			continue
		}

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
		if val != archiveCatalog {
			ret = append(ret, newCatalog.Lite())
		}
	}

	ret = append(ret, s.currentCatalog)

	return
}

func (s *Registry) archiveBlog() error {
	cmsClnt, cmsErr := s.getCMSClient()
	if cmsErr != nil {
		log.Printf("getCMSClient failed, err:%s", cmsErr.Error())
		return cmsErr
	}
	defer cmsClnt.Release()

	_, archives, _, commonErr := s.queryBlogCommon(cmsClnt)
	if commonErr != nil {
		log.Printf("queryBlogCommon failed, err:%s", commonErr.Error())
		return commonErr
	}

	var archiveCatalogPtr *cmsModel.CatalogLite
	preDuration := time.Duration(time.Now().UnixNano()) - time.Hour*24*2
	preTime := time.Unix(int64(preDuration.Seconds()), 0)
	archiveName := fmt.Sprintf("%04d年%02d月", preTime.Year(), preTime.Month())

	for _, val := range archives {
		if archiveName == val.Name {
			archiveCatalogPtr = val
			break
		}
	}

	if archiveCatalogPtr == nil {
		catalogPtr, catalogErr := cmsClnt.CreateCatalog(archiveName, "create archive catalog", s.archiveCatalog)
		if catalogErr == nil {
			return catalogErr
		}

		archiveCatalogPtr = catalogPtr.Lite()
	}

	archiveList, archiveErr := s.queryArticleList(cmsClnt, s.currentCatalog, nil)
	if archiveErr != nil {
		return archiveErr
	}

	for _, val := range archiveList {
		catalogs := []*cmsModel.CatalogLite{archiveCatalogPtr}
		for _, cv := range val.Catalog {
			if cv.Name == archiveCatalogPtr.Name {
				continue
			}

			if cv.ID != s.currentCatalog.ID {
				catalogs = append(catalogs, cv)
			}
		}

		val.Catalog = catalogs

		_, archiveErr := s.updateArticle(cmsClnt, val.ID, val.Title, val.Content, val.Catalog)
		if archiveErr != nil {
			return archiveErr
		}
	}

	return nil
}

// PostBlog post blog
func (s *Registry) PostBlog(res http.ResponseWriter, req *http.Request) {
	type postParam struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Catalog string `json:"catalog"`
	}

	type postResult struct {
		commonDef.Result
		Redirect string `json:"redirect"`
	}

	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := curSession.GetSessionInfo()
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

		cmsClient.BindSession(sessionInfo)

		catalogList, catalogErr := s.getCatalogs(param.Catalog, cmsClient)
		if catalogErr != nil {
			log.Printf("getCatalogs failed, err:%s", catalogErr.Error())
			result.ErrorCode = commonDef.Failed
			result.Reason = "提交Blog失败, 查询分类出错"
			break
		}

		if param.ID > 0 {
			_, err = s.updateArticle(cmsClient, param.ID, param.Title, param.Content, catalogList)
			if err != nil {
				result.ErrorCode = commonDef.Failed
				result.Reason = "提交Blog失败, 更新出错"
				break
			}
		} else {
			_, err = s.createArticle(cmsClient, param.Title, param.Content, catalogList)
			if err != nil {
				result.ErrorCode = commonDef.Failed
				result.Reason = "提交Blog失败, 保存出错"
				break
			}
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

// PostComment post comment
func (s *Registry) PostComment(res http.ResponseWriter, req *http.Request) {
	type postParam struct {
		Name    string `json:"name"`
		EMail   string `json:"email"`
		Message string `json:"message"`
		Origin  string `json:"origin"`
		Host    int    `json:"host"`
	}

	type postResult struct {
		commonDef.Result
		Redirect string `json:"redirect"`
	}

	result := &postResult{}
	for {
		param := &postParam{}
		err := net.ParseJSONBody(req, param)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		if param.Name == "" || param.EMail == "" || param.Origin == "" || param.Host == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数,输入参数为空"
			break
		}

		cmsClient, cmsErr := s.getCMSClient()
		if cmsErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "留言失败, 系统出错"
			break
		}
		defer cmsClient.Release()

		_, err = cmsClient.CreateComment(param.Message, param.Name, &cmsModel.Host{Code: param.Host, Type: cmsModel.ARTICLE}, 0)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "留言失败, 保存出错"
			break
		}

		result.ErrorCode = commonDef.Success
		result.Redirect = param.Origin
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

// ReplyComment reply comment
func (s *Registry) ReplyComment(res http.ResponseWriter, req *http.Request) {
	type postParam struct {
		Message string `json:"message"`
		Origin  string `json:"origin"`
		Host    int    `json:"host"`
	}

	type postResult struct {
		commonDef.Result
		Redirect string `json:"redirect"`
	}

	curSession := s.sessionRegistry.GetSession(res, req)
	sessionInfo := curSession.GetSessionInfo()
	result := &postResult{}
	for {
		param := &postParam{}
		err := net.ParseJSONBody(req, param)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		if param.Origin == "" || param.Host == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数,输入参数为空"
			break
		}

		cmsClient := cmsClient.NewClient(s.cmsService)
		defer cmsClient.Release()

		cmsClient.BindSession(sessionInfo)

		authPtr, _ := curSession.GetOption(commonCommon.AuthAccount)
		entityPtr := authPtr.(*casModel.Entity)

		_, err = cmsClient.CreateComment(param.Message, entityPtr.Name, &cmsModel.Host{Code: param.Host, Type: cmsModel.COMMENT}, 0)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "回复失败, 保存出错"
			break
		}

		result.ErrorCode = commonDef.Success
		result.Redirect = param.Origin
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

// DeleteComment delete comment
func (s *Registry) DeleteComment(res http.ResponseWriter, req *http.Request) {
	type postParam struct {
		Origin string `json:"origin"`
		Host   int    `json:"host"`
	}

	type postResult struct {
		commonDef.Result
		Redirect string `json:"redirect"`
	}

	curSession := s.sessionRegistry.GetSession(res, req)
	sessionInfo := curSession.GetSessionInfo()
	result := &postResult{}
	for {
		param := &postParam{}
		err := net.ParseJSONBody(req, param)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		if param.Origin == "" || param.Host == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数,输入参数为空"
			break
		}

		cmsClient := cmsClient.NewClient(s.cmsService)
		defer cmsClient.Release()

		cmsClient.BindSession(sessionInfo)

		commentList, _, commentErr := cmsClient.FilterComment(&cmsModel.Host{Code: param.Host, Type: cmsModel.COMMENT}, nil)
		if commentErr == nil && len(commentList) > 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "删除Comment失败,包含回复信息"
			break
		}

		_, err = cmsClient.DeleteComment(param.Host)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "删除Comment失败,删除数据出错"
			break
		}

		result.ErrorCode = commonDef.Success
		result.Redirect = param.Origin
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

// SettingBlog setting blog
func (s *Registry) SettingBlog(res http.ResponseWriter, req *http.Request) {
	type postResult struct {
		commonDef.Result
		Redirect string `json:"redirect"`
	}

	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := curSession.GetSessionInfo()
	result := &postResult{}
	for {
		param := &model.Setting{}
		err := net.ParseJSONBody(req, param)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		if param.Name == "" || param.Domain == "" || param.EMail == "" || param.ICP == "" {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数,输入参数为空"
			break
		}

		cmsClient := cmsClient.NewClient(s.cmsService)
		defer cmsClient.Release()

		cmsClient.BindSession(sessionInfo)

		catalogList := []*cmsModel.CatalogLite{s.systemCatalog}

		title := settingTitle
		content, contentErr := json.Marshal(param)
		if contentErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "保存Blog设置失败"
			break
		}

		if param.ID > 0 {
			_, err = s.updateArticle(cmsClient, param.ID, title, string(content), catalogList)
			if err != nil {
				result.ErrorCode = commonDef.Failed
				result.Reason = "保存Blog设置失败, 更新出错"
				break
			}
		} else {
			_, err = s.createArticle(cmsClient, title, string(content), catalogList)
			if err != nil {
				result.ErrorCode = commonDef.Failed
				result.Reason = "保存Blog设置失败, 保存出错"
				break
			}
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

func (s *Registry) createArticle(clnt cmsClient.Client, title, content string, catalogs []*cmsModel.CatalogLite) (ret *cmsModel.ArticleView, err error) {
	blogArticle, blogErr := clnt.CreateArticle(title, content, catalogs)
	if blogErr != nil {
		err = blogErr
		log.Printf("CreateArticle failed, err:%s", err.Error())
		return
	}

	ret = blogArticle
	return
}

func (s *Registry) deleteArticle(clnt cmsClient.Client, id int) (ret *cmsModel.ArticleView, err error) {
	blogArticle, blogErr := clnt.DeleteArticle(id)
	if blogErr != nil {
		err = blogErr
		log.Printf("DeleteArticle failed, err:%s", err.Error())
		return
	}

	ret = blogArticle
	return
}

func (s *Registry) updateArticle(clnt cmsClient.Client, id int, title, content string, catalogs []*cmsModel.CatalogLite) (ret *cmsModel.ArticleView, err error) {
	blogArticle, blogErr := clnt.UpdateArticle(id, title, content, catalogs)
	if blogErr != nil {
		err = blogErr
		log.Printf("UpdateArticle failed, err:%s", err.Error())
		return
	}

	ret = blogArticle
	return
}

func (s *Registry) deleteCatalog(clnt cmsClient.Client, id int) (ret *cmsModel.CatalogView, err error) {
	blogCatalog, blogErr := clnt.DeleteCatalog(id)
	if blogErr != nil {
		err = blogErr
		log.Printf("DeleteCatalog failed, err:%s", err.Error())
		return
	}

	ret = blogCatalog
	return
}

func (s *Registry) queryComments(clnt cmsClient.Client, id int, pageFilter *util.PageFilter) (ret []*model.CommentView, err error) {
	blogComment, _, blogErr := clnt.FilterComment(&cmsModel.Host{Code: id, Type: cmsModel.ARTICLE}, nil)
	if blogErr != nil {
		err = blogErr
		log.Printf("FilterComment failed, err:%s", err.Error())
		return
	}

	ret = []*model.CommentView{}
	for _, val := range blogComment {
		view := &model.CommentView{ID: val.ID, Content: val.Content, Creater: val.Creater, CreateDate: val.CreateDate}
		replyComment, _, replyErr := clnt.FilterComment(val.Host(), nil)

		view.Reply = []interface{}{}
		if replyErr == nil {
			for _, sv := range replyComment {
				rv := &model.CommentView{ID: sv.ID, Content: sv.Content, Creater: sv.Creater, CreateDate: sv.CreateDate}
				view.Reply = append(view.Reply, rv)
			}
		}

		ret = append(ret, view)
	}

	return
}
