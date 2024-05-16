CREATE TABLE IF NOT EXISTS users (
                                     id bigserial PRIMARY KEY,
                                     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
                                     name text NOT NULL,
                                     email text UNIQUE NOT NULL,
                                     password_hash bytea NOT NULL,
                                     activated bool NOT NULL,
                                     version integer NOT NULL DEFAULT 1
);
CREATE TABLE IF NOT EXISTS cars(
    id serial primary key,
    model varchar,
    brand varchar,
    year int,
    price float,
    color varchar,
    isUsed bool,
    userId serial references users(id),
    categoryName varchar references category(name)
);


