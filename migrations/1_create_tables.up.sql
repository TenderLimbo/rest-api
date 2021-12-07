CREATE TABLE IF NOT EXISTS books (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(100) NOT NULL UNIQUE,
                                     price NUMERIC(8) NOT NULL,
                                     genre INT NOT NULL,
                                     amount INT NOT NULL
);

CREATE TABLE IF NOT EXISTS genres (
                                      id SERIAL PRIMARY KEY,
                                      name VARCHAR(100) NOT NULL UNIQUE
);

INSERT INTO genres VALUES (1, 'adventure');
INSERT INTO genres VALUES (2, 'classics');
INSERT INTO genres VALUES (3, 'fantasy');