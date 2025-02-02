create table if not exists users
(
    id              SERIAL PRIMARY KEY,
    name            varchar not null,
    surname         varchar not null,
    patronymic      varchar,
    address         TEXT    not null,
    passport_number varchar not null
);

create table if not exists tasks
(
    id   SERIAL PRIMARY KEY,
    name varchar NOT NULL
);

create table if not exists user_tasks
(

    id         SERIAL PRIMARY KEY,
    user_id    INTEGER references users (id) ON DELETE CASCADE,
    task_id    INTEGER REFERENCES tasks (id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time   TIMESTAMP
)
