CREATE TABLE songs (
  id BIGSERIAL PRIMARY KEY,
  release_date date NOT NULL,
  link varchar(255),
  lyrics text
);
