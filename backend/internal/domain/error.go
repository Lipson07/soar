package domain

import "errors"

var (
	ErrNotFound             = errors.New("Ничего не найдена")
	ErrProjectNotFound      = errors.New("проект не найден")
	ErrUserNotFound         = errors.New("пользователь не найден")
	ErrTaskNotFound         = errors.New("задача не найдена")
	ErrUserAlreadyInProject = errors.New("пользоваеть уже есть в проекте")
	ErrEmailAlreadyExists   = errors.New("email уже используется")
	ErrInvalidID            = errors.New("некорректный ID")
	ErrInvalidCredentials   = errors.New("неверный email или пароль")
	ErrUnauthorized         = errors.New("не авторизован")
	ErrInternalServer       = errors.New("внутренняя ошибка сервера")
)