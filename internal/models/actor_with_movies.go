package models

type ActorWithMovies struct {
	ActorId       int     `json:"actor_id"`
	ActorName     string  `json:"actor_name"`
	ActorGender   string  `json:"actor_gender"`
	ActorBirthday string  `json:"actor_birthday"`
	Movies        []Movie `json:"movies"`
}

func NewActorWithMovies(actor Actor, movies []Movie) ActorWithMovies {
	return ActorWithMovies{
		ActorId:       actor.Id,
		ActorName:     actor.Name,
		ActorGender:   actor.Gender,
		ActorBirthday: actor.Birthday,
		Movies:        movies,
	}
}

/*
select a.id as "actor_id",
a.name as "actor_name",
a.gender as "actor_gender",
a.birthday as "actor_birthday",
m.id as "movie_id",
m.name as "movie_name",
m.description as "movie_description",
m.date as "movie_date",
m.rate as "movie_rate"
from actors a
join actors_to_movies am
on am.actor_id = a.id
join movies m
on m.id = am.movie_id;
 actor_id |     actor_name     | actor_gender | actor_birthday | movie_id | movie_name |      movie_description       | movie_date | movie_rate
----------+--------------------+--------------+----------------+----------+------------+------------------------------+------------+------------
        1 | Jody Foster        | female       | 1982-12-12     |        1 | Sumerki    | ffffffffffffffffffffffffff   | 2012-10-10 |         10
        1 | Jody Foster        | female       | 1982-12-12     |        2 | Kaptery    | Великое кино о великой войне | 2012-10-10 |         10
       23 | Skarlette Johanson | female       | 1985-12-12     |        2 | Kaptery    | Великое кино о великой войне | 2012-10-10 |         10
       24 | Valera ZHMA        | male         | 1970-12-12     |        2 | Kaptery    | Великое кино о великой войне | 2012-10-10 |         10
       25 | Alan Rickman       | male         | 1946-12-12     |        1 | Sumerki    | ffffffffffffffffffffffffff   | 2012-10-10 |         10

*/
