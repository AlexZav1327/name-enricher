-- +migrate Up
CREATE TABLE enriched_user (
    name VARCHAR NOT NULL PRIMARY KEY,
    age INT NOT NULL,
    gender VARCHAR NOT NULL,
    country VARCHAR NOT NULL
);