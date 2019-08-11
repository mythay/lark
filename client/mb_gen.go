package client

import (
	sqlx "github.com/jmoiron/sqlx"
)

type mbGen struct {
	db *sqlx.DB
}
type mbHost struct {
	Name     string
	Ipaddr   string
	Port     uint16
	Serial   string
	Baud     uint32
	Databits uint8
	Parity   string
	Stopbits uint8
	Interval uint16
}

func (gen *mbGen) getHosts() ([]mbHost, error) {
	hosts := []mbHost{}
	err := gen.db.Select(&hosts, `select [name], [ipaddr], [port], [serial], [baud], [databits], [parity], [stopbits], [interval] from [mb_host]`)
	return hosts, err
}
