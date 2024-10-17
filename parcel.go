package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, _ := s.db.Exec("INSERT INTO parcel (Client, Status, Address, created_at) VALUES (:Client, :Status, :Address, :created_at)",
		//sql.Named("Number", p.Number),
		sql.Named("Client", p.Client),
		sql.Named("Status", p.Status),
		sql.Named("Address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	// верните идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	return int(id), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	p := Parcel{}
	row := s.db.QueryRow("SELECT Number, Client, Status, Address, created_at FROM parcel WHERE Number = :Number", sql.Named("Number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	// заполните объект Parcel данными из таблицы
	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT Number, Client, Status, Address, created_at FROM parcel WHERE Client = :Client", sql.Named("Client", client))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		res = append(res, p)
	}
	// заполните срез Parcel данными из таблицы
	return res, err
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET Status = :Status WHERE Number = :Number",
		sql.Named("Status", status),
		sql.Named("Number", number))
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET Address = :Address WHERE Number = :Number, Status = :Status",
		sql.Named("Address", address),
		sql.Named("Number", number),
		sql.Named("Status", ParcelStatusRegistered))
	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE Number = :Number, Status = :Status",
		sql.Named("Number", number),
		sql.Named("Status", ParcelStatusRegistered))
	return err
}
