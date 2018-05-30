import React from 'react'
import PropTypes from 'prop-types'
import { Menu, Icon, Row } from 'antd'
import { Link } from 'dva/router'
import styles from './index.less'

const { SubMenu } = Menu

function MaintainHeader({ history, user, logoutHandler }) {
  const { location } = history

  const onMenuClick = ({ key }) => {
    if (key === 'logout') {
      logoutHandler()
    }
  }

  return (
    <div className={styles.maintain_content}>
      <Row type="flex" justify="end">
        <Menu
          selectedKeys={[location.pathname]}
          mode="horizontal"
          className={styles.menu}
          onClick={onMenuClick}
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
          <Menu.Item key="/maintain">
            <Link to="/maintain"><Icon type="tool" />Maintain</Link>
          </Menu.Item>
          <SubMenu
            key="/user"
            style={{ float: 'right' }}
            title={
              <span>
                <Icon type="user" />
                { user.name }
              </span>}
          >
            <Menu.Item key="logout">
              注销
            </Menu.Item>
          </SubMenu>
        </Menu>
      </Row>
    </div>
  )
}

MaintainHeader.propTypes = {
  history: PropTypes.object,
  user: PropTypes.object,
  logoutHandler: PropTypes.func,
}

export default MaintainHeader
