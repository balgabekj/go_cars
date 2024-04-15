CREATE TABLE IF NOT EXISTS owner(
    id serial primary key,
    name varchar,
    number varchar,
    email varchar,
    password varchar
);

CREATE TABLE IF NOT EXISTS car(
    id serial primary key,
    model varchar,
    brand varchar,
    year int,
    price float,
    color varchar,
    isUsed bool,
    ownerId serial references owner(id)
);


