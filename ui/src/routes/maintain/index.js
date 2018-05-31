import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row, Col } from 'antd'
import { SummaryTree } from '../common'

function MaintainPage({ maintain, loading, dispatch }) {
  const { summaryList } = maintain

  return (
    <Row type="flex" align="top">
      <Col md={4} lg={4} xl={4}>
        <SummaryTree summaryList={summaryList} />
      </Col>
      <Col md={20} lg={20} xl={20}>Right</Col>
    </Row>
  )
}

MaintainPage.propTypes = {
  loading: PropTypes.object,
  dispatch: PropTypes.func,
}

export default connect(({ maintain, loading, dispatch }) => ({ maintain, loading, dispatch }))(MaintainPage)
