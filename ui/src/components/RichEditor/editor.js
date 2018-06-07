import React, { Component } from 'react'
import PropTypes from 'prop-types'
import autobind from 'class-autobind'
import RichTextEditor from 'react-rte'
import { Tabs, Input } from 'antd'
import defaultFormat from './common'

const { TabPane } = Tabs
const { TextArea } = Input

export default class RichEditor extends Component {
  constructor(props) {
    super(props)

    autobind(this)

    const curFormat = defaultFormat
    let curValue = RichTextEditor.createEmptyValue()
    if (('value' in props) && props.value.length !== 0) {
      curValue = curValue.setContentFromString(props.value, curFormat)
    }

    this.state = {
      richValue: curValue,
      format: curFormat,
      placeholder: props.placeholder,
      editorStyle: props.editorStyle,
    }

    this.onChange = this.onChange.bind(this)
    this.onChangeSource = this.onChangeSource.bind(this)
  }

  componentWillReceiveProps(nextProps) {
    if (('value' in nextProps) && (this.props.value !== nextProps.value) && this.props.value.length === 0) {
      this.setState({ richValue: this.state.richValue.setContentFromString(nextProps.value, this.state.format) })
    }

    if (('placeholder' in nextProps) && (this.props.placeholder !== nextProps.placeholder)) {
      this.setState({ placeholder: nextProps.placeholder })
    }
  }

  onChange(value) {
    this.setState({ richValue: value })

    this.props.onChange(value.toString(this.state.format))
  }

  onChangeSource(event) {
    const source = event.target.value
    const oldValue = this.state.richValue

    this.setState({ richValue: oldValue.setContentFromString(source, this.state.format) })

    this.props.onChange(source)
  }

  render() {
    const { richValue, format, placeholder, editorStyle } = this.state

    return (
      <Tabs type="card">
        <TabPane tab="编辑器" key="richEditor">
          <div className="row">
            <RichTextEditor
              toolbarClassName="demo-toolbar"
              editorClassName="demo-editor"
              value={richValue}
              onChange={this.onChange}
              placeholder={placeholder}
              editorStyle={editorStyle}
            />
          </div>
        </TabPane>
        <TabPane tab="Markdown模式" key="markdown">
          <div className="row">
            <TextArea rows={27} cols={30} value={richValue.toString(format)} onChange={this.onChangeSource} />
          </div>
        </TabPane>
        <TabPane tab="预览" key="preView">
          <div className="row">
            <RichTextEditor value={richValue} editorStyle={editorStyle} readOnly />
          </div>
        </TabPane>
      </Tabs>
    )
  }
}

RichEditor.propTypes = {
  value: PropTypes.string,
  placeholder: PropTypes.string,
  editorStyle: PropTypes.object,
  onChange: PropTypes.func,
}
