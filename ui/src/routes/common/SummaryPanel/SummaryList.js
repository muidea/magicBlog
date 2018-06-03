import React from 'react'
import { Row, List, Col, Button } from 'antd'
import { Link } from 'dva/router'
import styles from './SummaryList.less'

function SummaryList({ summaryList, readOnly }) {
  const TitleText = ({ id, type, name }) => (
    <Link to={`/${type}/${id}`} ><h1>{name}</h1></Link>
  )
  const DescText = ({ description, creater, createDate }) => (
    <div>
      <div>{description}</div>
      <Row gutter={24} type="flex" align="middle">
        <Col xl={{ span: 18 }} md={{ span: 18 }}>
          <span>
            Post by {creater.name} on { createDate }
          </span>
        </Col>
        { !readOnly && <Col xl={{ span: 6 }} md={{ span: 6 }}>
          <Button className={styles.button}>编辑</Button>
          <Button className={styles.button}>删除</Button>
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
            title={<TitleText id={item.id} type={item.type} name={item.name} />}
            description={<DescText
              description={item.description}
              creater={item.creater}
              createDate={item.createDate}
            />}
          />
        </List.Item>
       )}
    />
  )
}

export default SummaryList
