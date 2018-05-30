import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row, Col } from 'antd'
import { LoginPanel } from '../../components'
import MaintainPanel from './maintain'

function MaintainPage({ app, maintain, loading, dispatch }) {
  const { isLogin } = app

  const onLogin = (value) => {
    dispatch({ type: 'maintain/loginUser', payload: { ...value } })
  }

  return (
    <Row type="flex">
      <Col span={24}>
        { !isLogin && <LoginPanel loading={loading} onLogin={onLogin} />}
        { isLogin && <MaintainPanel />}
      </Col>
    </Row>
  )
}

MaintainPage.propTypes = {
  maintain: PropTypes.object,
  loading: PropTypes.object,
  dispatch: PropTypes.func,
}

export default connect(({ app, maintain, loading, dispatch }) => ({ app, maintain, loading, dispatch }))(MaintainPage)
