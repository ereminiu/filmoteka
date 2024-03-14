CREATE TABLE movies (
    id serial primary key,
    name varchar(101) not null default '',
    description text not null default '',
    date DATE not null default '1997-07-14',
    rate int not null default(0)
);

CREATE TABLE actors (
    id serial primary key,
    name varchar(255) not null default '',
    gender varchar(255) not null default '',
    birthday varchar(255) not null default '1997-07-14',
    UNIQUE (name, birthday)
);

CREATE TABLE actors_to_movies (
    id serial primary key,
    actor_id int not null references actors,
    movie_id int not null references movies,
    unique (actor_id, movie_id)
);