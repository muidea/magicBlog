import React, { Component } from 'react'
import { Editor, EditorState } from 'draft-js'
import { stateFromMarkdown } from 'draft-js-import-markdown'
import PropTypes from 'prop-types'


class MarkdownEditor extends Component {
  constructor(props) {
    super(props)
    this.state = { editorState: EditorState.createWithContent(props.contentState) }
    this.onChange = editorState => this.setState({ editorState })
  }

  componentWillReceiveProps(nextProps) {
    if ('contentState' in nextProps) {
      this.setState({ editorState: EditorState.createWithContent(nextProps.contentState) })
    }
  }

  render() {
    return (
      <Editor editorState={this.state.editorState} onChange={this.onChange} readOnly />
    )
  }
}

const RichView = ({ value }) => {
  const contentState = stateFromMarkdown(value)
  return < MarkdownEditor contentState={contentState} />
}

RichView.propTypes = { value: PropTypes.string }

export default RichView
