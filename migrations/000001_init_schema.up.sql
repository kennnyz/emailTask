drop table if exists users;

create extension if not exists "uuid-ossp";

create table if not exists users (
                                     email VARCHAR(255),
                                     created_at TIMESTAMP DEFAULT NOW(),
                                     unique_code UUID DEFAULT uuid_generate_v4()
)
