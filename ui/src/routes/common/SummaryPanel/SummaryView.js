import React from 'react'
import PropTypes from 'prop-types'
import { Row, List, Col, Button } from 'antd'
import { Link } from 'dva/router'
import styles from './SummaryView.less'

function SummaryView({ summaryList, readOnly, onSelect, onModify, onDelete }) {
  const TitleText = ({ item }) => (
    <div>
      <a><h1 onClick={() => onSelect(item)}>{item.name}</h1></a>
    </div>
  )

  const DescText = ({ item }) => (
    <div>
      <div>{ item.description}</div>
      <Row gutter={24} type="flex" align="middle">
        <Col xl={{ span: 18 }} md={{ span: 18 }}>
          <span>
            Post by {item.creater.name} on { item.createDate }
          </span>
        </Col>
        { !readOnly && <Col xl={{ span: 6 }} md={{ span: 6 }}>
          <Button className={styles.button} onClick={() => onModify(item)} >编辑</Button>
          <Button className={styles.button} onClick={() => onDelete(item)} >删除</Button>
        </Col>
        }
      </Row>
    </div>
  )

  const MoreInfo = () => (
    <Row type="flex" justify="end">
      <Col><Link to="/catalog">More</Link></Col>
    </Row>
    )

  return (
    <List
      itemLayout="horizontal"
      dataSource={summaryList}
      footer={<MoreInfo />}
      renderItem={item => (
        <List.Item>
          <List.Item.Meta
            title={<TitleText item={item} />}
            description={<DescText item={item} />}
          />
        </List.Item>
       )}
    />
  )
}

SummaryView.propTypes = {
  summaryList: PropTypes.array,
  readOnly: PropTypes.bool,
  onSelect: PropTypes.func,
  onModify: PropTypes.func,
  onDelete: PropTypes.func,
}

export default SummaryView
