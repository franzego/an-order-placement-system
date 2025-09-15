create table products (
    id bigserial primary key,
    product_name varchar(255) not null,
    product_description text,
    stock numeric(10,2) not null,
    price numeric(10, 2) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp
);