CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(101) NOT NULL DEFAULT '',
    description text NOT NULL DEFAULT '',
    date DATE NOT NULL DEFAULT '1997-07-14',
    rate INT NOT NULL DEFAULT(0)
);

CREATE TABLE IF NOT EXISTS actors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT '',
    gender VARCHAR(255) NOT NULL DEFAULT '',
    birthday VARCHAR(255) NOT NULL DEFAULT '1997-07-14',
    UNIQUE (name, birthday)
);

CREATE TABLE IF NOT EXISTS actors_to_movies (
    id SERIAL PRIMARY KEY,
    actor_id INT NOT NULL references actors,
    movie_id INT NOT NULL references movies,
    UNIQUE (actor_id, movie_id)
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    uesrname VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);