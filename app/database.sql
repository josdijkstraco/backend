
CREATE SCHEMA jpar;

DROP TABLE jpar.user;
CREATE TABLE jpar.user (
	id integer GENERATED ALWAYS AS IDENTITY,
	username text,
	first_name text,
	last_name text,
	email text,
	password text,
	is_staff boolean,
	is_active boolean,
	is_superuser boolean
);

DROP TABLE jpar.review;
CREATE TABLE jpar.review (
	id integer GENERATED ALWAYS AS IDENTITY,
	user_id integer,
	product_id integer,
	name text,
	rating integer,
	comment text,
	created_at time
);

DROP TABLE jpar.product;
CREATE TABLE jpar.product (
	id integer GENERATED ALWAYS AS IDENTITY,
	user_id integer,
	name text,
	image text,
	brand text,
	category text,
	description text,
	rating integer,
	number_reviews integer,
	price decimal,
	count_in_stock integer,
	created_at time
);

DROP TABLE jpar.order;
CREATE TABLE jpar.order (
	id integer GENERATED ALWAYS AS IDENTITY,
	user_id integer,
	payment_method text,
	tax_price decimal,
	shipping_price decimal,
	total_price decimal,
	is_paid boolean,
	paid_at time,
	is_delivered boolean,
	delivered_at time,
	created_at time
);

DROP TABLE jpar.order_item;
CREATE TABLE jpar.order_item (
	id integer GENERATED ALWAYS AS IDENTITY,
	order_id integer,
	product text,
	name text,
	quantity integer,
	price decimal,
	image text
);

DROP TABLE jpar.shipping_address;
CREATE TABLE jpar.shipping_address (
	id integer GENERATED ALWAYS AS IDENTITY,
	order_id integer,
	address text,
	city text,
	postal_code text,
	country text,
	shipping_price decimal
);


