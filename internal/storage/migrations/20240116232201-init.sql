-- +migrate Up
CREATE TABLE enriched_user (
    name VARCHAR NOT NULL,
    surname VARCHAR NOT NULL,
    patronymic VARCHAR,
    age INT NOT NULL,
    gender VARCHAR NOT NULL,
    country VARCHAR NOT NULL
);