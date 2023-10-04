drop table if exists invoice_item;
drop table if exists invoice;
drop table if exists product;
drop table if exists customer;

create table customer (
    id varchar(50) primary key
    ,name varchar(50) not null
);

create table product (
    id varchar(50) primary key
    ,name varchar(50) not null
    ,price money not null
);

create table invoice (
    id varchar(50) primary key
    ,customer_id varchar(50) not null
    ,date date not null
    ,foreign key (customer_id) references customer(id)
);

create table invoice_item (
    id varchar(50) primary key
    ,invoice_id varchar(50) not null
    ,product_id varchar(50) not null
    ,quantity int not null
    ,foreign key (invoice_id) references invoice(id) on delete cascade
    ,foreign key (product_id) references product(id)
);

begin transaction;
insert into customer (id, name) values ('1', 'John');
insert into customer (id, name) values ('2', 'Mary');
insert into customer (id, name) values ('3', 'Peter');

insert into product (id, name, price) values ('1', 'Apple', 1.99);
insert into product (id, name, price) values ('2', 'Orange', 2.99);
insert into product (id, name, price) values ('3', 'Banana', 3.99);

insert into invoice (id, customer_id, date) values ('1', '1', '2019-01-01');
insert into invoice (id, customer_id, date) values ('2', '2', '2019-01-02');
insert into invoice (id, customer_id, date) values ('3', '3', '2019-01-03');

insert into invoice_item (id, invoice_id, product_id, quantity) values ('1', '1', '1', 1);
insert into invoice_item (id, invoice_id, product_id, quantity) values ('2', '1', '2', 2);
insert into invoice_item (id, invoice_id, product_id, quantity) values ('3', '2', '1', 3);
insert into invoice_item (id, invoice_id, product_id, quantity) values ('4', '2', '2', 4);
insert into invoice_item (id, invoice_id, product_id, quantity) values ('5', '2', '3', 5);
insert into invoice_item (id, invoice_id, product_id, quantity) values ('6', '3', '1', 6);
insert into invoice_item (id, invoice_id, product_id, quantity) values ('7', '3', '2', 7);

commit transaction;
```

