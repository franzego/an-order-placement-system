create table store (
    id bigserial primary key,
    product_name varchar(255) not null,
    product_description text,
    available_quantity numeric(10,2) not null,
    reserved_quantity numeric(10,2) not null,
    price numeric(10, 2) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);