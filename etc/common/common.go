package common

import "errors"

var (
	ErrInternal = errors.New("Внутренная ошибка сервера. Пожалуйста попробуйте позднее или свяжитесь с системным администратором")
	ErrNotFound = errors.New("Запись не найдена")
)
