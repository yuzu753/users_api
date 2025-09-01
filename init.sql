-- Sample tenant tables for development
CREATE TABLE "users_tenant1" (
    id                    SERIAL                       PRIMARY KEY,
    user_name             TEXT                         UNIQUE NOT NULL,
    email                 TEXT                         UNIQUE NOT NULL,
    password_hash         TEXT                         NOT NULL,
    otp_secret_key        TEXT,
    type                  INTEGER                      NOT NULL,
    authority_data        TEXT,
    failed_count          INTEGER                      NOT NULL DEFAULT 0,
    unlock_at             TIMESTAMP WITHOUT TIME ZONE,
    is_reset_password     BOOLEAN                      DEFAULT FALSE,
    last_update_password  TIMESTAMP WITHOUT TIME ZONE  NOT NULL DEFAULT now()
);

CREATE TABLE "users_tenant2" (
    id                    SERIAL                       PRIMARY KEY,
    user_name             TEXT                         UNIQUE NOT NULL,
    email                 TEXT                         UNIQUE NOT NULL,
    password_hash         TEXT                         NOT NULL,
    otp_secret_key        TEXT,
    type                  INTEGER                      NOT NULL,
    authority_data        TEXT,
    failed_count          INTEGER                      NOT NULL DEFAULT 0,
    unlock_at             TIMESTAMP WITHOUT TIME ZONE,
    is_reset_password     BOOLEAN                      DEFAULT FALSE,
    last_update_password  TIMESTAMP WITHOUT TIME ZONE  NOT NULL DEFAULT now()
);

-- Insert sample data for testing
INSERT INTO "users_tenant1" (user_name, email, password_hash, type, last_update_password) VALUES
('john_doe', 'john@example.com', '$2a$10$hashedpassword1', 1, now()),
('jane_smith', 'jane@example.com', '$2a$10$hashedpassword2', 2, now()),
('admin_user', 'admin@example.com', '$2a$10$hashedpassword3', 0, now());

INSERT INTO "users_tenant2" (user_name, email, password_hash, type, last_update_password) VALUES
('alice_wonder', 'alice@tenant2.com', '$2a$10$hashedpassword4', 1, now()),
('bob_builder', 'bob@tenant2.com', '$2a$10$hashedpassword5', 2, now());