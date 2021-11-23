CREATE TABLE IF NOT EXISTS books (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(100) NOT NULL,
                                     price NUMERIC(2) NOT NULL,
                                     genre INT NOT NULL,
                                     amount INT NOT NULL
);

CREATE TABLE IF NOT EXISTS genres (
                                      id SERIAL PRIMARY KEY,
                                      name VARCHAR(100) NOT NULL
);

INSERT INTO genres VALUES (1, 'adventure');
INSERT INTO genres VALUES (2, 'classics');
INSERT INTO genres VALUES (3, 'fantasy');