CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20)
);

CREATE TABLE cakes (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    price INT NOT NULL,
    weight INT NOT NULL,
    full_description TEXT NOT NULL,
    active BOOLEAN DEFAULT TRUE
);

CREATE TABLE delivery_points (
    id SERIAL PRIMARY KEY,
    address TEXT NOT NULL,
    rating INT CHECK (rating >= 0 AND rating <= 10),
    working_hours VARCHAR(50) NOT NULL,
    contact_phone VARCHAR(20)
);

CREATE TABLE deliveries (
    id SERIAL PRIMARY KEY,
    point_id INT REFERENCES delivery_points(id) ON DELETE CASCADE,
    cost INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    weight INT NOT NULL
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_status VARCHAR(50) NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    payment_method VARCHAR(50) NOT NULL
);

CREATE TABLE order_cakes (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    cake_id INT REFERENCES cakes(id) ON DELETE CASCADE
);

CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    level VARCHAR(10),
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);