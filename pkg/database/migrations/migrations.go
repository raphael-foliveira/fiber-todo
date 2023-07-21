package database


var Schema = `
CREATE TABLE IF NOT EXISTS todo (
    id SERIAL PRIMARY KEY,
    title VARCHAR,
    description VARCHAR,
    completed BOOLEAN
);
`

var Migrations = []string{
	``,
}