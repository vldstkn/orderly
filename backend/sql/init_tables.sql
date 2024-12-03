CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT,
    number_completed_orders INT DEFAULT 0
);

CREATE TABLE provider(
    id INT PRIMARY KEY REFERENCES users(id),
    rating FLOAT DEFAULT 0.0
);

CREATE TABLE customer(
    id INT PRIMARY KEY REFERENCES users(id),
    points int
);

CREATE TABLE category(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE products(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    provider_id INT REFERENCES provider(id) NOT NULL,
    quantity INT DEFAULT 0,
    rating FLOAT DEFAULT 0.0,
    cat_id INT REFERENCES category(id),
    discount INT DEFAULT 0,
    price INT NOT NULL
);

CREATE TYPE order_status as ENUM('Pending', 'Processing', 'Failed', 'Canceled', 'Completed');

CREATE TABLE orders(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    customer_id INT NOT NULL REFERENCES customer(id),
    provider_id INT NOT NULL REFERENCES provider(id),
    product_id INT NOT NULL REFERENCES products(id),
    status order_status DEFAULT 'Pending',
    total_price INT NOT NULL
);

CREATE TABLE reviews(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    product_id INT NOT NULL REFERENCES products(id),
    customer_id INT NOT NULL REFERENCES customer(id),
    description TEXT,
    rating INT CHECK (rating >=0 AND rating <= 5)
)