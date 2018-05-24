import React from 'react'
import PropTypes from 'prop-types'
import { Row, Col } from 'antd'
import LoginForm from './form'
import styles from './index.less'

const LoginPanel = ({ loading, onLogin }) => (
  <Row>
    <Col md={2} lg={4}><LoginForm loading={loading} onLogin={onLogin} /></Col>
    <Col md={{ offset: 6 }} lg={{ offset: 8 }}><div className={styles.info}>
      <h3>账号登陆</h3>
      <p>请输入正确的账号信息进行登录</p>
      <p>登录成功后允许进行文章编写和发布</p>
    </div>
    </Col>
  </Row>
  )

LoginPanel.propTypes = {
  loading: PropTypes.object,
  onLogin: PropTypes.func,
}

export default LoginPanel
