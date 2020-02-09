package handler

import (
	"fmt"
	"time"

	"github.com/muidea/magicBatis/client"
	"github.com/muidea/magicBatis/def"
	"github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicBlog/model"	
)

// CommonHandler common handler
type CommonHandler interface {
	UpdateSysConfig(cfg *model.SysConfig) (*model.SysConfig, error)
	GetSysConfig() (*model.SysConfig, error)

	QueryOpLog(pageFilter *util.PageFilter) ([]*model.OpLog, int64, error)
	WriteOpLog(account, address, memo string) (*model.OpLog, error)	
}

// NewCommonHandler 新建CommonHandler
func NewCommonHandler(clnt client.Client) (CommonHandler, error) {
	handler := &commonHandler{batisClnt: clnt}

	return handler, nil
}

type commonHandler struct {
	batisClnt client.Client
}

func (s *commonHandler) UpdateSysConfig(cfg *model.SysConfig) (ret *model.SysConfig, err error) {
	err = s.batisClnt.UpdateEntity(cfg)
	if err == nil {
		ret = cfg
	}

	return
}

func (s *commonHandler) GetSysConfig() (ret *model.SysConfig, err error) {
	cfg := &model.SysConfig{}
	err = s.batisClnt.QueryEntity(cfg)
	if err == nil {
		ret = cfg
	}

	return
}


func (s *commonHandler) QueryOpLog(pageFilter *util.PageFilter) (ret []*model.OpLog, total int64, err error) {
	ret = []*model.OpLog{}
	queryfilter := &def.QueryFilter{}
	if pageFilter != nil {
		queryfilter.Page(pageFilter)
	}
	queryfilter.SortFilter = &util.SortFilter{AscSort: false, FieldName: "CreateTime"}

	total, err = s.batisClnt.BatchQueryEntity(&ret, queryfilter)

	return
}

func (s *commonHandler) WriteOpLog(account, address, memo string) (ret *model.OpLog, err error) {
	if account == "" {
		err = fmt.Errorf("操作人员信息不能为空")
		return
	}
	oplog := &model.OpLog{Account: account, Address: address, Memo: memo, CreateTime: time.Now().Format("2006-01-02 15:04:05")}
	err = s.batisClnt.InsertEntity(oplog)
	if err == nil {
		ret = oplog
	}

	return
}
