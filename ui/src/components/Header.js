import React from 'react'
import { Menu, Icon } from 'antd'
import { Link } from 'dva/router'

function Header({ location }) {
  return (
    <Menu
      selectedKeys={[location.pathname]}
      mode="horizontal"
      theme="dark"
    >
      <Menu.Item key="/catalog">
        <Link to="/catalog"><Icon type="bars" />Catalog</Link>
      </Menu.Item>
      <Menu.Item key="/content">
        <Link to="/content"><Icon type="home" />Content</Link>
      </Menu.Item>
      <Menu.Item key="/404">
        <Link to="/page-you-dont-know"><Icon type="frown-circle" />404</Link>
      </Menu.Item>
      <Menu.Item key="/antd">
        <a href="https://github.com/dvajs/dva" target="_new">dva</a>
      </Menu.Item>
    </Menu>
  )
}

export default Header