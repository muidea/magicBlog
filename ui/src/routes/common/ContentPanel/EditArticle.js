import React from 'react'
import PropTypes from 'prop-types'
import { Button, Form, Input } from 'antd'
import { RichEditor } from '../../../components'
import styles from './index.less'

const FormItem = Form.Item

const EditArticle = ({
  onSubmit,
  contentItem,
  form: {
    getFieldDecorator,
    validateFieldsAndScroll,
  },
}) => {
  const handleOk = () => {
    validateFieldsAndScroll((errors, values) => {
      if (errors) {
        return
      }

      const { parent } = contentItem
      const { id, type } = parent

      values = {
        ...values,
        parent: { id, type },
      }

      onSubmit(values)
    })
  }

  return (
    <div className={styles.form}>
      <div>
        <form>
          <FormItem label="标题" hasFeedback>
            {getFieldDecorator('title', {
              initialValue: contentItem.title,
              rules: [{ required: true }],
            })(<Input size="large" onPressEnter={handleOk} placeholder="输入标题" />)}
          </FormItem>
          <FormItem label="内容">
            {getFieldDecorator('content', {
              initialValue: contentItem.content,
            })(<RichEditor />)}
          </FormItem>
        </form>
      </div>
      <div>
        <Button type="primary" size="large" onClick={handleOk}>
            添加
        </Button>
      </div>
    </div>
  )
}

EditArticle.propTypes = {
  form: PropTypes.object,
  onSubmit: PropTypes.func,
  contentItem: PropTypes.object,
}

export default Form.create()(EditArticle)
