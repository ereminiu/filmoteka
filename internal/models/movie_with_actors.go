package models

type MovieWithActors struct {
	MovieId          int    `db:"movie_id"`
	MovieName        string `db:"movie_name"`
	MovieDescription string `db:"movie_description"`
	MovieDate        string `db:"movie_date"`
	MovieRate        int    `db:"movie_rate"`
	ActorId          int    `db:"actor_id"`
	ActorName        string `db:"actor_name"`
	ActorGender      string `db:"actor_gender"`
	ActorBirthday    string `db:"actor_birthday"`
}

/*

movie_id |
movie_name |
movie_description |
movie_date |
movie_rate |
actor_id |
actor_name     |
actor_gender |
actor_birthday
*/

/*

select m.id as "movie_id", m.name as "movie_name", m.description as "movie_description", m.date as "movie_date", m.rate as "movie_rate", a.id as "actor_id", a.name as "actor_name", a.gender as "actor_gender", a.birthdate as "actor_birthday"
from movies m
left join actor_to_movie am
on am.movie_id = m.id
left join actors a
on a.id = am.actor_id;

*/
