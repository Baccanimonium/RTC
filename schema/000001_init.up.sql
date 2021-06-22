CREATE TABLE users
(
    id                 serial       not null unique,
    name               varchar(255) not null,
    login              varchar(255) not null unique,
    password           varchar(255) not null,
    about              varchar(255),
    address            varchar(255),
    phone              varchar(255)
);

CREATE TABLE doctor (
    id                 serial not null unique,
    id_user            int references users (id) on delete cascade not null,
    salary             money,
    qualifications     varchar(255),
    contacts           varchar(255)
);

CREATE TABLE patient (
    id                 serial not null unique,
    id_user            int references users (id) on delete cascade not null,
    description        varchar(500),
    recovered          boolean      not null default false
);

CREATE TABLE course (
    id                 serial not null unique,
    id_patient         int    references patient (id) on delete cascade not null,
    id_users           int    references users (id) on delete set null,
    title              varchar(100),
    description        varchar(500),
    time_start         varchar(10),
    time_end           varchar(10)
);

CREATE TABLE schedule (
    id                 serial not null unique,
    id_user            int    references users (id) on delete cascade not null,
    title              varchar(100),
    description        varchar(500)
);

CREATE TABLE consultation (
    id                 serial not null unique,
    id_patient         int    references patient (id) on delete set null,
    id_users           int    references users (id) on delete cascade not null,
    id_course          int    references course (id) on delete set null,
    title              varchar(100),
    time_start         varchar(16),
    time_end           varchar(16)
);

CREATE TABLE event (
    id                 serial not null unique,
    time_start         varchar(16),
    time_end           varchar(16),
    title              varchar(100),
    id_course          int    references course (id) on delete cascade not null,
    type               varchar(100),
    description        varchar(500),
    accepted           boolean      not null default false
);

CREATE TABLE schedule_event
(
    id                 serial         not null unique,
    id_schedule        int references schedule (id) on delete cascade not null,
    id_event           int references event (id) on delete cascade not null
);