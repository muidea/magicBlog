import React from 'react'
import PropTypes from 'prop-types'
import { Tree, Icon } from 'antd'

const { TreeNode } = Tree

function SummaryTree({ summaryList }) {
  const getNodeIcon = (type) => {
    if (type === 'catalog') {
      return <Icon type="tags-0" />
    } else if (type === 'link') {
      return <Icon type="share-alt" />
    } else if (type === 'media') {
      return <Icon type="picture" />
    } else {
      return <Icon type="file" />
    }
  }

  const renderTreeNodes = (data) => {
    return data.map((item) => {
      if (item.subSummary) {
        return (
          <TreeNode icon={<Icon type="tags-o" />} title={item.name} key={`/${item.type}/${item.id}`} dataRef={item}>
            {renderTreeNodes(item.subSummary)}
          </TreeNode>
        )
      }

      return <TreeNode icon={getNodeIcon(item.type)} title={item.name} key={`/${item.type}/${item.id}`} dataRef={item} />
    })
  }

  const onSelect = (selectedKeys, info) => {
    console.log('selected', selectedKeys, info)
  }

  return (
    <Tree
      showIcon
      onSelect={onSelect}
    >
      { renderTreeNodes(summaryList) }
    </Tree>
  )
}

SummaryTree.propTypes = {
  summaryList: PropTypes.array,
}

export default SummaryTree
