package models

type MovieWithActors struct {
	MovieId          int     `json:"movie_id"`
	MovieName        string  `json:"movie_name"`
	MovieDescription string  `json:"movie_description"`
	MovieDate        string  `json:"movie_date"`
	MovieRate        int     `json:"movie_rate"`
	Actors           []Actor `json:"actors"`
}

func NewMovieWithActors(movie Movie, actors []Actor) MovieWithActors {
	return MovieWithActors{
		MovieId:          movie.Id,
		MovieName:        movie.Name,
		MovieDescription: movie.Description,
		MovieDate:        movie.Date,
		MovieRate:        movie.Rate,
		Actors:           actors,
	}
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
