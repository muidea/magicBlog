import React from 'react'
import { Row, List, Col } from 'antd'
import { Link } from 'dva/router'

function SummaryList({ summaryList }) {
  const TitleText = ({ id, type, name }) => (
    <Link to={`/${type}/${id}`} ><h1>{name}</h1></Link>
  )
  const DescText = ({ description, creater, createDate }) => (
    <div>
      <div>{description}</div>
      <span>
         Post by {creater.name} on { createDate }
      </span>
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
