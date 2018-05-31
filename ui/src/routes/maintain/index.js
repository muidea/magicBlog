import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row, Col } from 'antd'
import { SummaryTree, SummaryContent } from '../common'

function MaintainPage({ maintain, dispatch }) {
  const { summaryList, currentSelect } = maintain

  const onSelect = (value) => {
    dispatch({ type: 'maintain/querySelectContent', payload: { ...value } })
  }

  return (
    <Row type="flex" align="top">
      <Col md={4} lg={4} xl={4}>
        <SummaryTree summaryList={summaryList} onSelect={onSelect} />
      </Col>
      <Col md={20} lg={20} xl={20}>
        <SummaryContent contentData={currentSelect} />
      </Col>
    </Row>
  )
}

MaintainPage.propTypes = {
  dispatch: PropTypes.func,
}

export default connect(({ maintain, dispatch }) => ({ maintain, dispatch }))(MaintainPage)
