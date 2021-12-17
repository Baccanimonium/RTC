CREATE TABLE users
(
    id       serial not null unique,
    name     varchar(40),
    surname  varchar(40),
    avatar   varchar(40),
    login    varchar(255) not null unique,
    password varchar(255) not null,
    about    varchar(255),
    address  varchar(255),
    phone    varchar(255),
    promoted boolean,
    deleted  boolean
);

CREATE TABLE doctor
(
    id             serial not null unique,
    id_user        int references users (id) on delete cascade not null unique,
    salary         float8,
    qualifications varchar(255),
    contacts       varchar(255)
);

CREATE TABLE patient
(
    id                serial not null unique,
    id_user           int references users (id) on delete cascade not null unique,
    id_current_doctor int references users (id) on delete set null,
    description       varchar(500),
    recovered         boolean not null default false
);

CREATE TABLE patients_log
(
    id                serial not null unique,
    id_patient        int references patient (id) on delete cascade not null unique,
    text              text,
    file              varchar(40)
    created_at        bigint
    log_type          varchar(25)
);

CREATE TABLE consultation
(
    id            serial not null unique,
    id_patient    int references patient (id) on delete set null,
    id_user       int references users (id) on delete cascade not null,
    time_start    varchar(16),
    last          int,
    offline       boolean,
    doctor_joined varchar(25)
);

CREATE TABLE event
(
    id                    serial not null unique,
    id_patient            int references patient (id) on delete cascade not null,
    id_doctor             int references users (id) on delete cascade not null,
    created_at            varchar(25),
    last_till             varchar(10),
    at_days               jsonb,
    title                 varchar(100),
    description           varchar(500),
    notify_doctor         boolean,
    remind_in_advance     int default 0,
    confirmation_time     int default 3600,
    weight                int default 1,
    requires_confirmation boolean not null default false
);

CREATE TABLE roles
(
    id          serial not null unique,
    title       varchar(25) not null unique,
)

CREATE TABLE user_role
(
    id          serial not null unique,
    id_user     int references users (id) on delete cascade not null,
    id_event    int references roles (id) on delete cascade not null
);

CREATE TABLE tasks
(
    id         serial not null unique,
    id_patient int references patient (id) on delete cascade not null,
    id_user    int references users (id) on delete cascade not null,
    id_event   int references event (id) on delete cascade not null,
    weight     int default 0,
);

CREATE TABLE consultation_files
(
    id              serial not null unique,
    notes           text not null,
    id_consultation int references users (consultation) on delete cascade not null,
    for_doctor      boolean default true,
    files           varchar(40)[],
);