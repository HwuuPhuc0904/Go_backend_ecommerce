-- Create the PostgreSQL database and user
CREATE DATABASE ecommerce_db;
CREATE USER ecommerce_user WITH ENCRYPTED PASSWORD 'secure_password';

-- Grant all privileges on the database to the user
GRANT ALL PRIVILEGES ON DATABASE ecommerce_db TO ecommerce_user;

-- You can add additional initialization SQL commands here, such as creating tables or inserting initial data.