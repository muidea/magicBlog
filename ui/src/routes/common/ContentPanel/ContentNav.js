import React from 'react'
import PropTypes from 'prop-types'
import { Tree, Icon } from 'antd'

const { TreeNode } = Tree

function ContentNav({ itemList, onSelect }) {
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
      if (item.subItem) {
        return (
          <TreeNode icon={<Icon type="tags-o" />} title={item.name} key={`/${item.type}/${item.id}`} dataRef={item}>
            {renderTreeNodes(item.subItem)}
          </TreeNode>
        )
      }

      return <TreeNode icon={getNodeIcon(item.type)} title={item.name} key={`/${item.type}/${item.id}`} dataRef={item} />
    })
  }

  const onSelectItem = (selectedKeys, { selectedNodes }) => {
    if (selectedKeys.length > 0) {
      onSelect(selectedNodes[0].props.dataRef)
    }
  }

  return (
    <Tree
      showIcon
      onSelect={onSelectItem}
    >
      { renderTreeNodes(itemList) }
    </Tree>
  )
}

ContentNav.propTypes = {
  itemList: PropTypes.array,
  onSelect: PropTypes.func,
}

export default ContentNav
