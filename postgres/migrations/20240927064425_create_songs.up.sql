CREATE TABLE song (
  id BIGSERIAL PRIMARY KEY,
  releaseDate date NOT NULL,
  link varchar(255),
  lyrics text
);
