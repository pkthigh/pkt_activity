package storage

// MYSQL mysql db
type MYSQL string

const ()

func (mysql MYSQL) String() string {
	return string(mysql)
}

// AREA redis area
type AREA int

const ()

// Int to int
func (area AREA) Int() int {
	return int(area)
}

// COLL mongo coll
type COLL string

const ()

// Int to int
func (coll COLL) String() string {
	return string(coll)
}
