sqldimel
========

SQLDimel is an SQL DML query builder for Golang.

It just generates the INSERT, UPDATE and DELETE queries to be passed to the underlining database.

Using a processor, it is possible to customize the way parameters are generated on the SQL (with ?, $1, etc).

Processors using "?" and "$1, $2" are included.

NOTE: when using WHERE, ALWAYS use "?", this is the default placeholder that will be replaced by the processor.

INSTALLATION
------------

	go get github.com/RangelReale/sqldimel

	import "github.com/RangelReale/sqldimel"


USAGE
-----

	b := sqldimel.NewBuilder("user")
	b.Add("id", 1).
		Add("name", "Monte Marto").
		Add("dob", time.Now()).
		Add("optional", nil).
		Add("weight", 80.2).
		Where("id = ? and weight > ?", 1, 70.2)

	// db = *sql.DB
	res, err = b.Exec(db, sqldimel.UPDATE)
	// or
	// res, err = db.Exec(b.Output(sqldimel.UPDATE), b.OutputParams(sqldimel.UPDATE))
