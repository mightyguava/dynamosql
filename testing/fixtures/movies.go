package fixtures

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/require"
)

// Movies is a fixture with a movies table and a few top movies.
var Movies = Fixture{
	Table: *moviesTable,
	Create: &dynamodb.CreateTableInput{
		TableName: moviesTable,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("title"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("year"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("title"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("year"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	},
	Data: func(t *testing.T, client *dynamodb.DynamoDB) {
		var movies []Movie
		err := json.Unmarshal([]byte(movieData), &movies)
		require.NoError(t, err)
		for _, m := range movies {
			_, err := client.PutItem(&dynamodb.PutItemInput{
				TableName: moviesTable,
				Item: map[string]*dynamodb.AttributeValue{
					"title": {S: aws.String(m.Title)},
					"year":  {N: aws.String(strconv.Itoa(m.Year))},
					"info": {
						M: map[string]*dynamodb.AttributeValue{
							"actors": {
								SS: m.Info.Actors,
							},
							"directors": {
								SS: m.Info.Directors,
							},
							"genres": {
								SS: m.Info.Genres,
							},
							"imageURL": {
								S: aws.String(m.Info.ImageURL),
							},
							"plot": {
								S: aws.String(m.Info.Plot),
							},
							"rank": {
								N: aws.String(strconv.Itoa(m.Info.Rank)),
							},
							"rating": {
								N: aws.String(strconv.FormatFloat(m.Info.Rating, 'g', -1, 64)),
							},
							"release_date": {
								S: aws.String(m.Info.ReleaseDate),
							},
							"running_time_secs": {
								S: aws.String(strconv.Itoa(m.Info.RunningTimeSecs)),
							},
						},
					},
				},
			})
			require.NoError(t, err)
		}
	},
}

// moviesTable is the name of the movies table in the Movies
var moviesTable = aws.String("movies")

// Movie is a container for movie data
type Movie struct {
	Title string    `json:"title"`
	Year  int       `json:"year"`
	Info  MovieInfo `json:"info"`
}

type MovieInfo struct {
	Actors          []*string `json:"actors"`
	Directors       []*string `json:"directors"`
	Genres          []*string `json:"genres"`
	ImageURL        string    `json:"image_url"`
	Plot            string    `json:"plot"`
	Rank            int       `json:"rank"`
	Rating          float64   `json:"rating"`
	ReleaseDate     string    `json:"release_date"`
	RunningTimeSecs int       `json:"running_time_secs"`
}

const movieData = `[
  {
    "year": 2013,
    "title": "Prisoners",
    "info": {
      "directors": [
        "Denis Villeneuve"
      ],
      "release_date": "2013-08-30T00:00:00Z",
      "rating": 8.2,
      "genres": [
        "Crime",
        "Drama",
        "Thriller"
      ],
      "image_url": "http://ia.media-imdb.com/images/M/MV5BMTg0NTIzMjQ1NV5BMl5BanBnXkFtZTcwNDc3MzM5OQ@@._V1_SX400_.jpg",
      "plot": "When Keller Dover's daughter and her friend go missing, he takes matters into his own hands as the police pursue multiple leads and the pressure mounts. But just how far will this desperate father go to protect his family?",
      "rank": 3,
      "running_time_secs": 9180,
      "actors": [
        "Hugh Jackman",
        "Jake Gyllenhaal",
        "Viola Davis"
      ]
    }
  },
  {
    "year": 2013,
    "title": "The Hunger Games: Catching Fire",
    "info": {
      "directors": [
        "Francis Lawrence"
      ],
      "release_date": "2013-11-11T00:00:00Z",
      "genres": [
        "Action",
        "Adventure",
        "Sci-Fi",
        "Thriller"
      ],
      "image_url": "http://ia.media-imdb.com/images/M/MV5BMTAyMjQ3OTAxMzNeQTJeQWpwZ15BbWU4MDU0NzA1MzAx._V1_SX400_.jpg",
      "plot": "Katniss Everdeen and Peeta Mellark become targets of the Capitol after their victory in the 74th Hunger Games sparks a rebellion in the Districts of Panem.",
      "rank": 4,
      "running_time_secs": 8760,
      "actors": [
        "Jennifer Lawrence",
        "Josh Hutcherson",
        "Liam Hemsworth"
      ]
    }
  },
  {
    "year": 2013,
    "title": "Thor: The Dark World",
    "info": {
      "directors": [
        "Alan Taylor"
      ],
      "release_date": "2013-10-30T00:00:00Z",
      "genres": [
        "Action",
        "Adventure",
        "Fantasy"
      ],
      "image_url": "http://ia.media-imdb.com/images/M/MV5BMTQyNzAwOTUxOF5BMl5BanBnXkFtZTcwMTE0OTc5OQ@@._V1_SX400_.jpg",
      "plot": "Faced with an enemy that even Odin and Asgard cannot withstand, Thor must embark on his most perilous and personal journey yet, one that will reunite him with Jane Foster and force him to sacrifice everything to save us all.",
      "rank": 5,
      "actors": [
        "Chris Hemsworth",
        "Natalie Portman",
        "Tom Hiddleston"
      ]
    }
  },
  {
    "year": 2013,
    "title": "This Is the End",
    "info": {
      "directors": [
        "Evan Goldberg",
        "Seth Rogen"
      ],
      "release_date": "2013-06-03T00:00:00Z",
      "rating": 7.2,
      "genres": [
        "Comedy",
        "Fantasy"
      ],
      "image_url": "http://ia.media-imdb.com/images/M/MV5BMTQxODE3NjM1Ml5BMl5BanBnXkFtZTcwMzkzNjc4OA@@._V1_SX400_.jpg",
      "plot": "While attending a party at James Franco's house, Seth Rogen, Jay Baruchel and many other celebrities are faced with the apocalypse.",
      "rank": 6,
      "running_time_secs": 6420,
      "actors": [
        "James Franco",
        "Jonah Hill",
        "Seth Rogen"
      ]
    }
  },
  {
    "year": 2013,
    "title": "Insidious: Chapter 2",
    "info": {
      "directors": [
        "James Wan"
      ],
      "release_date": "2013-09-13T00:00:00Z",
      "rating": 7.1,
      "genres": [
        "Horror",
        "Thriller"
      ],
      "image_url": "http://ia.media-imdb.com/images/M/MV5BMTg0OTA5ODIxNF5BMl5BanBnXkFtZTcwNTUzNDg4OQ@@._V1_SX400_.jpg",
      "plot": "The haunted Lambert family seeks to uncover the mysterious childhood secret that has left them dangerously connected to the spirit world.",
      "rank": 7,
      "running_time_secs": 6360,
      "actors": [
        "Patrick Wilson",
        "Rose Byrne",
        "Barbara Hershey"
      ]
    }
  },
  {
    "year": 2013,
    "title": "World War Z",
    "info": {
      "directors": [
        "Marc Forster"
      ],
      "release_date": "2013-06-02T00:00:00Z",
      "rating": 7.1,
      "genres": [
        "Action",
        "Adventure",
        "Horror",
        "Sci-Fi",
        "Thriller"
      ],
      "image_url": "http://ia.media-imdb.com/images/M/MV5BMTg0NTgxMjIxOF5BMl5BanBnXkFtZTcwMDM0MDY1OQ@@._V1_SX400_.jpg",
      "plot": "United Nations employee Gerry Lane traverses the world in a race against time to stop the Zombie pandemic that is toppling armies and governments, and threatening to destroy humanity itself.",
      "rank": 8,
      "running_time_secs": 6960,
      "actors": [
        "Brad Pitt",
        "Mireille Enos",
        "Daniella Kertesz"
      ]
    }
  },
  {
    "year": 2014,
    "title": "X-Men: Days of Future Past",
    "info": {
      "directors": [
        "Bryan Singer"
      ],
      "release_date": "2014-05-21T00:00:00Z",
      "genres": [
        "Action",
        "Adventure",
        "Fantasy",
        "Sci-Fi"
      ],
      "image_url": "http://ia.media-imdb.com/images/M/MV5BMTQ0NzIwNTA1MV5BMl5BanBnXkFtZTgwNjY2OTcwMDE@._V1_SX400_.jpg",
      "plot": "The X-Men send Wolverine to the past to change a major historical event that could globally impact man and mutant kind.",
      "rank": 9,
      "actors": [
        "Jennifer Lawrence",
        "Hugh Jackman",
        "Michael Fassbender"
      ]
    }
  },
  {
    "year": 2014,
    "title": "Transformers: Age of Extinction",
    "info": {
      "directors": [
        "Michael Bay"
      ],
      "release_date": "2014-06-25T00:00:00Z",
      "genres": [
        "Action",
        "Adventure",
        "Sci-Fi"
      ],
      "image_url": "http://ia.media-imdb.com/images/M/MV5BMTQyMDA5Nzg0Nl5BMl5BanBnXkFtZTgwNzA4NDcxMDE@._V1_SX400_.jpg",
      "plot": "A mechanic and his daughter make a discovery that brings down Autobots and Decepticons - and a paranoid government official - on them.",
      "rank": 10,
      "actors": [
        "Mark Wahlberg",
        "Nicola Peltz",
        "Jack Reynor"
      ]
    }
  },
  {
    "year": 2013,
    "title": "Now You See Me",
    "info": {
      "directors": [
        "Louis Leterrier"
      ],
      "release_date": "2013-05-21T00:00:00Z",
      "rating": 7.3,
      "genres": [
        "Crime",
        "Mystery",
        "Thriller"
      ],
      "image_url": "http://ia.media-imdb.com/images/M/MV5BMTY0NDY3MDMxN15BMl5BanBnXkFtZTcwOTM5NzMzOQ@@._V1_SX400_.jpg",
      "plot": "An FBI agent and an Interpol detective track a team of illusionists who pull off bank heists during their performances and reward their audiences with the money.",
      "rank": 11,
      "running_time_secs": 6900,
      "actors": [
        "Jesse Eisenberg",
        "Common",
        "Mark Ruffalo"
      ]
    }
  }
]
`
