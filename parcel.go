package main

import (
	"database/sql"
	"fmt"
	"log"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (Number, Client, Status, Address, CreatedAt) VALUES (:Number, :Client, :Status, :Address, :CreatedAt)",
		sql.Named("Number", p.Number),
		sql.Named("Client", p.Client),
		sql.Named("Status", p.Status),
		sql.Named("Address", p.Address),
		sql.Named("CreatedAt", p.CreatedAt))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	fmt.Println(res.LastInsertId())
	return 0, nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	p := Parcel{}
	row := s.db.QueryRow("SELECT Number, Client, Status, Addres, CreatedAt FROM parcel WHERE Number = :Number", sql.Named("Number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return p, err
	}
	// заполните объект Parcel данными из таблицы
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT Number, Client, Status, Addres, CreatedAt FROM parcel WHERE Client = :Client", sql.Named("Client", client))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	// заполните срез Parcel данными из таблицы
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET Status = :Status WHERE Nunber = :Number",
		sql.Named("Status", status),
		sql.Named("Number", number))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	p := Parcel{}
	row := s.db.QueryRow("SELECT Status FROM parcel WHERE Number = :Number", sql.Named("Number", number))
	err := row.Scan(&p.Status)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if fmt.Sprintln(p) == "registered" {
		_, err := s.db.Exec("UPDATE parcel SET Addres = :Addres WHERE Nunber = :Number",
			sql.Named("Addres", address),
			sql.Named("Number", number))
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	p := Parcel{}
	row := s.db.QueryRow("SELECT Status FROM parcel WHERE Number = :Number", sql.Named("Number", number))
	err := row.Scan(&p.Status)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if fmt.Sprintln(p) == "registered" {
		_, err = s.db.Exec("DELETE FROM parcel WHERE Number = :Number", sql.Named("Number", number))
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
