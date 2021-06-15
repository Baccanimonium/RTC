CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    login         varchar(255) not null unique,
    password      varchar(255) not null,
    about         varchar(255),
    address       varchar(255),
    phone         varchar(255),
);

CREATE TABLE doctors (

);


CREATE TABLE patients (

);

CREATE TABLE group (

);

