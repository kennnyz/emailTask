drop table if exists users;

create extension if not exists "uuid-ossp";

create table if not exists users (
                                     email VARCHAR(255),
                                     unique_code UUID DEFAULT uuid_generate_v4()
)
