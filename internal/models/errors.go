package models

import "errors"

var (
	ErrNoRecord           = errors.New("没有匹配记录")
	ErrInvalidCredentials = errors.New("认证凭据无效")
	ErrDuplicateEmail     = errors.New("邮箱已被使用")
)
