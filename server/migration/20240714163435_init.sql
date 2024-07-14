-- +goose Up
-- +goose StatementBegin
create table gauge
(
    name  text not null
        constraint gauge_name_pk
            primary key,
    value double precision
);

create table counter
(
    name  text not null
        constraint counter_name_pk
            primary key,
    value bigint
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
