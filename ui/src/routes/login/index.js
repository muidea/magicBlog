import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row, Col } from 'antd'
import { LoginPanel } from '../common'

function LoginPage({ loading, dispatch }) {
  const onLogin = (value) => {
    dispatch({ type: 'app/loginUser', payload: { ...value } })
  }

  return (
    <Row type="flex">
      <Col span={24}>
        <LoginPanel loading={loading} onLogin={onLogin} />
      </Col>
    </Row>
  )
}

LoginPage.propTypes = {
  loading: PropTypes.object,
  dispatch: PropTypes.func,
}

export default connect(({ loading, dispatch }) => ({ loading, dispatch }))(LoginPage)
