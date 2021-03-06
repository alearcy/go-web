DROP TABLE IF EXISTS users;


DROP TABLE IF EXISTS sessions;


CREATE TABLE sessions ( id serial PRIMARY KEY,
                                          uuid VARCHAR(64) NOT NULL UNIQUE,
                                                                    user_id INTEGER NOT NULL UNIQUE,
                                                                                             created_at TIMESTAMP NOT NULL);


CREATE TABLE users ( id SERIAL PRIMARY KEY,
                                       name VARCHAR(255) NOT NULL,
                                                         surname VARCHAR(255) NOT NULL,
                                                                              email VARCHAR(255) NOT NULL,
                                                                                                 password VARCHAR(255) NOT NULL,
                                                                                                                       role INT NOT NULL DEFAULT 1,
                                                                                                                                                 created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                                                                                                                                                                       updated_at TIMESTAMP NOT NULL DEFAULT NOW());