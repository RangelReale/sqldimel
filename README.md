sqldimel
========

SQLDimel is an SQL DML query builder for Golang.

It just generates the INSERT, UPDATE and DELETE queries to be passed to the underlining database.

Using a processor, it is possible to customize the way parameters are generated on the SQL (with ?, $1, etc).
