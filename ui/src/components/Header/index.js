import React from 'react'
import { Menu, Icon, Row } from 'antd'
import { Link } from 'dva/router'
import styles from './index.less'

function Header({ history }) {
  const { location } = history

  return (
    <div className={styles.content}>
      <Row type="flex" justify="end">
        <Menu
          selectedKeys={[location.pathname]}
          mode="horizontal"
          className={styles.menu}
        >
          <Menu.Item key="/">
            <Link to="/"><Icon type="home" />Home</Link >
          </Menu.Item>
          <Menu.Item key="/catalog">
            <Link to="/catalog"><Icon type="appstore-o" />Post</Link>
          </Menu.Item>
          <Menu.Item key="/contact">
            <Link to="/contact"><Icon type="profile" />Contact</Link>
          </Menu.Item>
          <Menu.Item key="/about">
            <Link to="/about"><Icon type="idcard" />About</Link>
          </Menu.Item>
        </Menu>
      </Row>
      <Row type="flex" justify="space-around" align="middle">
        <div className={styles.info}>
          <h1>Muidea Blog</h1>
          <span>写作也是一种生活</span>
        </div>
      </Row>
    </div>
  )
}

export default Header
