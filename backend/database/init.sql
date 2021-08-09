CREATE TABLE IF NOT EXISTS fermantationBarrelTemperture (
    ID          SERIAL PRIMARY KEY,
    temperture  DOUBLE PRECISION,
    contentID   INTEGER,
    createdOn   TIMESTAMP DEFAULT now(),
);