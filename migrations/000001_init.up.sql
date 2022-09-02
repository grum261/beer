CREATE TABLE users (
    id int PRIMARY KEY generated always AS identity,
    username varchar(50) UNIQUE NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    avatar varchar(255) UNIQUE,
    bio varchar(500),
    password_hash varchar(255) NOT NULL
);

CREATE TABLE users_friends_requests (
    user_sender_id int NOT NULL REFERENCES users,
    user_reciever_id int NOT NULL REFERENCES users,
    sended_at timestamp NOT NULL DEFAULT NOW(),
    request_status int NOT NULL CHECK (request_status IN (0, 1, 2, 3))
);

CREATE UNIQUE INDEX uq_idx_users_id ON users_friends_requests (user_sender_id, user_reciever_id);

CREATE INDEX ON users_friends_requests (request_status);

CREATE TABLE breweries (
    id int PRIMARY KEY generated always AS identity,
    name varchar(255) UNIQUE NOT NULL,
    description varchar(1000),
    founded_at timestamp NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

-- возможно стоит добавить дату основания пивоварни в индекс
CREATE TABLE beers (
    id int PRIMARY KEY generated always AS identity,
    name varchar(255) NOT NULL UNIQUE,
    description varchar(500),
    image varchar(255) NOT NULL UNIQUE,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE brewery_beers (
    brewery_id int NOT NULL REFERENCES breweries,
    beer_id int NOT NULL REFERENCES beers
);

CREATE UNIQUE INDEX uq_idx_brewery_id_beer_id ON brewery_beers (beer_id, brewery_id);

CREATE TABLE tags (
    id int PRIMARY KEY generated always AS identity,
    name varchar(100) NOT NULL UNIQUE,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

CREATE TABLE posts (
    id int PRIMARY KEY generated always AS identity,
    score int NOT NULL CHECK (
        score <= 10
        AND score > 0
    ),
    review_text varchar(2000),
    image varchar(255) NOT NULL UNIQUE,
    user_id int NOT NULL REFERENCES users,
    beer_id int NOT NULL REFERENCES beers,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE INDEX ON posts (user_id, beer_id, score);