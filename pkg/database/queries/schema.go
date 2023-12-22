package queries

var Schema = `
    CREATE TABLE IF NOT EXISTS todo (
        id SERIAL PRIMARY KEY,
        title VARCHAR UNIQUE,
        description VARCHAR,
        completed BOOLEAN
    );
`
