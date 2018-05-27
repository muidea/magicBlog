import React from 'react'
import { Link } from 'dva/router'
import { Divider, BackTop } from 'antd'
import { config } from 'utils'
import styles from './Footer.less'

const { footerText } = config

const Footer = () => (
  <div className={styles.footer}>
    <Divider /><BackTop />
    <div>
      { footerText } <Divider type="vertical" /> <Link to="/maintain">管理</Link>
    </div>
  </div>)

export default Footer
