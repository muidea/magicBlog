import React from 'react'
import PropTypes from 'prop-types'
import { Button, Row, Form, Input } from 'antd'
import styles from './index.less'

const FormItem = Form.Item

const LoginForm = ({
  loading,
  onLogin,
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
      onLogin(values)
    })
  }

  return (
    <div className={styles.form}>
      <form>
        <FormItem hasFeedback>
          {getFieldDecorator('account', {
            rules: [
              {
                required: true,
                message: '请输入账号',
              },
            ],
          })(<Input size="large" onPressEnter={handleOk} placeholder="账号" />)}
        </FormItem>
        <FormItem hasFeedback>
          {getFieldDecorator('password', {
            rules: [
              {
                required: true,
                message: '请输入密码',
              },
            ],
          })(<Input size="large" type="password" onPressEnter={handleOk} placeholder="密码" />)}
        </FormItem>
        <Row>
          <Button type="primary" size="large" onClick={handleOk} loading={loading.effects.login}>
            登陆
          </Button>
        </Row>
      </form>
    </div>
  )
}

LoginForm.propTypes = {
  form: PropTypes.object,
  loading: PropTypes.object,
  onLogin: PropTypes.func,
}

export default Form.create()(LoginForm)
