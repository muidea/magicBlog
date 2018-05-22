import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Link } from 'dva/router'
import { Row, List, Col } from 'antd'
import styles from './index.css'

function IndexPage({ index }) {
  const { summaryList } = index

  const DescText = ({ creater, createDate }) => (
    <span>
      Post by {creater.name} on { createDate }
    </span>
  )

  const MoreInfo = () => (
    <Row type="flex" justify="end">
      <Col><Link to="/catalog">More</Link></Col>
    </Row>
  )

  return (
    <Row className={styles.normal}>
      <Col span={16} offset={4}>
        <List
          itemLayout="horizontal"
          dataSource={summaryList}
          footer={<MoreInfo />}
          renderItem={item => (
            <List.Item>
              <List.Item.Meta
                title={<a href="/contact"><h1>{item.name}</h1></a>}
                description={<DescText creater={item.creater} createDate={item.createDate} />}
              />
            </List.Item>
          )}
        />
      </Col>
    </Row>
  )
}

IndexPage.propTypes = {
  index: PropTypes.object,
}

export default connect(({ index }) => ({ index }))(IndexPage)
