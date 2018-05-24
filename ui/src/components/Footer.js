import React from 'react'
import { Link } from 'dva/router'
import { Divider, Icon, BackTop } from 'antd'
import styles from './Footer.less'

const Footer = () => (
  <div className={styles.footer}>
    <Divider /><BackTop />
    <div>
      <Icon type="copyright" /> 2018 muidea.com <Divider type="vertical" /> <Link to="/maintain">管理</Link>
    </div>
  </div>)

export default Footer
