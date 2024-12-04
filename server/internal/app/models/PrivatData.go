package models

import "time"

type PrivateData struct {
	ID         int       `db:"id"`          // Уникальный идентификатор данных
	UserID     int       `db:"user_id"`     // Внешний ключ на пользователя
	DataType   string    `db:"data_type"`   // Тип данных (пароль, заметка, бинарные данные, карта и т.д.)
	DataKey    string    `db:"data_key"`    // Ключ данных (например, имя сайта, тип карты)
	TextData   string    `db:"text_data"`   // Текстовые данные
	BinaryData []byte    `db:"binary_data"` // Бинарные данные (например, файл)
	CardNumber string    `db:"card_number"` // Номер банковской карты (если это данные карты)
	CreatedAt  time.Time `db:"created_at"`  // Дата создания
	UpdatedAt  time.Time `db:"updated_at"`  // Дата последнего обновления
	MetaInfo   string    `db:"meta_info"`   // Метаинформация (в формате JSON)
}
