CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS profiles
(
    user_id INT NOT NULL,
    first_name VARCHAR(125) NOT NULL,
    last_name VARCHAR(125) NOT NULL,
    phone_number VARCHAR(20),
    profile_pic VARCHAR(255),

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);