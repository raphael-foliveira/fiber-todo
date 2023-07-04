create table if not exists todo (
    id serial primary key,
    title varchar,
    description varchar,
    completed boolean
);