create table rating_place
(
    id         bigserial,
    place_id   bigint                  not null,
    rating     integer,
    review     varchar(250),
    created_at timestamp default now() not null,
    user_id    integer                 not null
);

alter table rating_place
    owner to postgres;

