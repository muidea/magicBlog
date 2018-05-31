import React from 'react'
import PropTypes from 'prop-types'
import { Row, Col } from 'antd'
import { MaintainHeader } from './Header'
import Footer from './Footer'
import styles from './MaintainLayout.less'

function MaintainLayout({ history, user, logoutHandler, children }) {
  return (
    <div>
      <div className={styles.header}>
        <MaintainHeader history={history} user={user} logoutHandler={logoutHandler} />
      </div>
      <Row className={styles.content}>
        <Col span={16} offset={4}>
          {children}
        </Col>
      </Row>
      <div className={styles.footer}><Footer /></div>
    </div>
  )
}

MaintainLayout.propTypes = {
  history: PropTypes.object,
  user: PropTypes.object,
  logoutHandler: PropTypes.func,
  children: PropTypes.object,
}


export default MaintainLayout
