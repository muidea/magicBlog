import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row, Col } from 'antd'
import { SummaryTree, ContentPanel } from '../common'
import styles from './index.less'

function MaintainPage({ maintain, dispatch }) {
  const { summaryList, currentSelect } = maintain

  const onSelect = (value) => {
    dispatch({ type: 'maintain/querySelectContent', payload: { ...value } })
  }

  const onAddCatalog = () => {

  }

  const onAddArticle = () => {

  }

  return (
    <Row type="flex" align="top">
      <Col md={4} lg={4} xl={4} className={styles.nav}>
        <SummaryTree summaryList={summaryList} onSelect={onSelect} />
      </Col>
      <Col md={20} lg={20} xl={20}>
        <ContentPanel contentData={currentSelect} onAddCatalog={onAddCatalog} onAddArticle={onAddArticle} />
      </Col>
    </Row>
  )
}

MaintainPage.propTypes = {
  dispatch: PropTypes.func,
}

export default connect(({ maintain, dispatch }) => ({ maintain, dispatch }))(MaintainPage)
