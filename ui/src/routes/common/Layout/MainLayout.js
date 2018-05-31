import React from 'react'
import { Row, Col } from 'antd'
import { NormalHeader } from './Header'
import Footer from './Footer'
import styles from './MainLayout.less'

function MainLayout({ history, children }) {
  return (
    <div>
      <div className={styles.header}>
        <NormalHeader history={history} />
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

export default MainLayout
