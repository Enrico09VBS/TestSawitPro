/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE users (
	id serial PRIMARY KEY,
	full_name VARCHAR ( 100 ) NOT NULL,
  password VARCHAR (100) NOT NULL,
  phone_number ( 20 ) UNIQUE NOT NULL,
  login_count int DEFAULT 0
);

INSERT INTO test (name) VALUES ('test1');
INSERT INTO test (name) VALUES ('test2');
