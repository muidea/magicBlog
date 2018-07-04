import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Tag, Input, Tooltip, Icon } from 'antd'

export default class EditableTagGroup extends Component {
  constructor(props) {
    super(props)
    if ('value' in props) {
      const { value } = props
      if (value) {
        this.state = { ...this.state, value }
      }
    }
  }

  state = {
    value: [],
    inputVisible: false,
    inputValue: '',
  }

  componentWillReceiveProps(nextProps) {
    if ('value' in nextProps) {
      const { value } = nextProps
      if (value) {
        this.setState({ value })
      }
    }
  }

  handleClose = (removedTag) => {
    const value = this.state.value.filter(tag => tag.name !== removedTag)
    this.setState({ value })

    const { onChange } = this.props
    if (onChange) {
      onChange(value)
    }
  }

  showInput = () => {
    this.setState({ inputVisible: true }, () => this.input.focus())
  }

  handleInputChange = (e) => {
    this.setState({ inputValue: e.target.value })
  }

  handleInputConfirm = () => {
    const { inputValue } = this.state
    let { value } = this.state
    let exitFlag = false
    for (const item of value) {
      if (item.name === inputValue) {
        exitFlag = true
        break
      }
    }
    if (inputValue && !exitFlag) {
      value = [...value, { name: inputValue, id: -1 }]
    }

    this.setState({
      value,
      inputVisible: false,
      inputValue: '',
    })

    const { onChange } = this.props
    if (onChange) {
      onChange(value)
    }
  }

  saveInputRef = (input) => {
    this.input = input
  }

  render() {
    const { value, inputVisible, inputValue } = this.state
    return (
      <div>
        {value && value.map((tag) => {
          const isLongTag = tag.name.length > 20
          const tagElem = (
            <Tag key={tag.name} closable={!this.props.readOnly} afterClose={() => this.handleClose(tag.name)}>
              {isLongTag ? `${tag.name.slice(0, 20)}...` : tag.name}
            </Tag>
          )
          return isLongTag ? <Tooltip title={tag.name} key={tag.name}>{tagElem}</Tooltip> : tagElem
        })}
        {inputVisible && !this.props.readOnly && (
          <Input
            ref={this.saveInputRef}
            type="text"
            size="small"
            style={{ width: 78 }}
            value={inputValue}
            onChange={this.handleInputChange}
            onBlur={this.handleInputConfirm}
            onPressEnter={this.handleInputConfirm}
          />
        )}
        {!inputVisible && !this.props.disableInput && !this.props.readOnly && (
          <Tag
            onClick={this.showInput}
            style={{ background: '#fff', borderStyle: 'dashed' }}
          >
            <Icon type="plus" /> 添加
          </Tag>
        )}
      </div>
    )
  }
}

EditableTagGroup.propTypes = {
  value: PropTypes.array,
  disableInput: PropTypes.bool,
  readOnly: PropTypes.bool,
  onChange: PropTypes.func,
}
