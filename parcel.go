package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// добавление строки с таблицу parcel
func (s ParcelStore) Add(p Parcel) (int, error) {
	// добавляем строку в таблицу
	res, err := s.db.Exec("INSERT INTO parcel (client,status,address,created_at) VALUES(:client,:status,:address,:created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	// получаем id последней записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// чтение строки из БД по id посылки
func (s ParcelStore) Get(number int) (Parcel, error) {
	// запрашиваем строку из таблицы
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number",
		sql.Named("number", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

// чтение строк из БД по ФИО клиента
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// срез Parcel
	var res []Parcel

	// запрашиваем строки из БД
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		return res, err
	}
	defer rows.Close()

	// распаковываем данные из БД
	for rows.Next() {
		var item Parcel

		err := rows.Scan(&item.Number, &item.Client, &item.Status, &item.Address, &item.CreatedAt)
		if err != nil {
			return res, err
		}

		res = append(res, item)
	}

	return res, nil
}

// смена статуса в БД у клиента по id
func (s ParcelStore) SetStatus(number int, status string) error {
	// запрос на изменение статуса в БД
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

// смена адреса в БД у клиента по id
func (s ParcelStore) SetAddress(number int, address string) error {

	// запрос на изменение адреса в БД
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	return nil
}

// удаление строки из БД по id
func (s ParcelStore) Delete(number int) error {

	// запрос на удаление строки из БД
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	return nil
}
