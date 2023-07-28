package migrations

var Schema = `
CREATE TABLE IF NOT EXISTS todo (
    id SERIAL PRIMARY KEY,
    title VARCHAR UNIQUE,
    description VARCHAR,
    completed BOOLEAN
);
`

var Migrations = []string{}
