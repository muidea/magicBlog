import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row, Col } from 'antd'
import { SummaryTree, ViewContent, EditCatalog, EditArticle } from '../common'
import styles from './index.less'

function MaintainPage({ maintain, dispatch }) {
  const { summaryList, action } = maintain

  const onSelect = (value) => {
    dispatch({ type: 'maintain/querySelectContent', payload: { ...value } })
  }

  const onAddCatalog = (parent) => {
    dispatch({ type: 'maintain/addCatalog', payload: { ...parent } })
  }

  const onSubmitCatalog = (value) => {
    dispatch({ type: 'maintain/submitCatalog', payload: { ...value } })
  }

  const onAddArticle = (parent) => {
    dispatch({ type: 'maintain/addArticle', payload: { ...parent } })
  }

  const onSubmitArticle = (value) => {
    dispatch({ type: 'maintain/submitArticle', payload: { ...value } })
  }

  const getContentPanel = () => {
    if (action.type === 'viewContent') {
      return <ViewContent contentData={action.value} onAddCatalog={onAddCatalog} onAddArticle={onAddArticle} />
    } else if (action.type === 'addCatalog') {
      return <EditCatalog contentItem={action.value} onSubmit={onSubmitCatalog} />
    } else if (action.type === 'addArticle') {
      return <EditArticle contentItem={action.value} onSubmit={onSubmitArticle} />
    } else {
      return <div>aaa</div>
    }
  }

  return (
    <Row type="flex" align="top">
      <Col md={6} lg={6} xl={6} className={styles.nav}>
        <SummaryTree summaryList={summaryList} onSelect={onSelect} />
      </Col>
      <Col md={18} lg={18} xl={18}>
        {getContentPanel()}
      </Col>
    </Row>
  )
}

MaintainPage.propTypes = {
  dispatch: PropTypes.func,
}

export default connect(({ maintain, dispatch }) => ({ maintain, dispatch }))(MaintainPage)
