import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row, Col } from 'antd'
import { LoginPanel } from '../../components'

function MaintainPage({ maintain, loading, dispatch }) {
  const { isLogin } = maintain

  const onLogin = (value) => {
    console.log(value)
  }

  return (
    <Row type="flex">
      <Col span={24}>
        { !isLogin && <LoginPanel loading={loading} onLogin={onLogin} />}
      </Col>
    </Row>
  )
}

MaintainPage.propTypes = {
  maintain: PropTypes.object,
  loading: PropTypes.object,
  dispatch: PropTypes.func,
}

export default connect(({ maintain, loading, dispatch }) => ({ maintain, loading, dispatch }))(MaintainPage)
