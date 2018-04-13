package core

import "muidea.com/magicCommon/model"

// Agent Center访问代理
type Agent interface {
	Start(account, password string) bool
	Stop()
	QuerySummary() []model.SummaryView
	QueryCatalog() []model.CatalogDetailView
	QuerySubCatalog() []model.CatalogDetailView
	QueryContent(id int) model.ArticleDetailView
}
