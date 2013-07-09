sqldimel
========

SQLDimel is an SQL DML query builder for Golang.

It just generates the INSERT, UPDATE and DELETE queries to be passed to the underlining database.

Using a processor, it is possible to customize the way parameters are generated on the SQL (with ?, $1, etc).

Processors using "?" and "$1, $2" are included.

INSTALLATION
------------

go get github.com/RangelReale/sqldimel

import "github.com/RangelReale/sqldimel"


USAGE
-----

	b := sqldimel.NewBuilder("user")
	b.Add("id", 1)
	b.Add("name", "Monte Marto")
	b.Add("dob", time.Now())
	b.Add("optional", nil)
	b.Add("weight", 80.2)

	b.Where("id = ? and weight > ?", 1, 70.2)

	// db = *sql.DB
	res, err = db.Exec(b.Output(sqldimel.UPDATE), b.OutputParams(sqldimel.UPDATE))
