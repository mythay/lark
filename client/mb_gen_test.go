package client

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func Test_mbGen_getModbusHosts(t *testing.T) {
	assert := assert.New(t)

	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return
	}
	defer db.Close()

	db.MustExec(`CREATE TABLE [main].[mb_host](
		[id] INTEGER PRIMARY KEY AUTOINCREMENT, 
		[name] TEXT, 
		[ipaddr] TEXT, 
		[port] INTEGER DEFAULT 502, 
		[serial] TEXT, 
		[baud] INTEGER, 
		[databits] INTEGER, 
		[parity] CHAR(1), 
		[stopbits] INTEGER, 
		[interval] INTEGER, 
		CHECK([parity] IN ('N', 'E', 'O')));
	  
	  /* Table data [mb_host] Record count: 2 */
	  INSERT INTO [mb_host]([id], [name], [ipaddr], [port], [serial], [baud], [databits], [parity], [stopbits], [interval]) VALUES(1, 'test', '127.0.0.1', 502, '', 0, 0, 'N', 0, 5);
	  INSERT INTO [mb_host]([id], [name], [ipaddr], [port], [serial], [baud], [databits], [parity], [stopbits], [interval]) VALUES(2, 'test2', '', 0, 'COM16', 115200, 8, 'E', 0, 10);
	  `)

	mb := mbGen{db}

	hosts, err := mb.getHosts()

	assert.Nil(err)
	assert.Equal(2, len(hosts))
}
