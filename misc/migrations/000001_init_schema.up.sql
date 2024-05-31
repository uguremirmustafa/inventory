CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    avatar TEXT,
    active_group_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS item (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    user_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    item_type_id BIGINT NOT NULL,
    manufacturer_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS item_type (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS item_info (
    id BIGSERIAL PRIMARY KEY,
    item_id BIGINT NOT NULL,
    purchase_date TIMESTAMP,
    expiration_date TIMESTAMP,
    last_used TIMESTAMP,
    purchase_location VARCHAR(255),
    price BIGINT,
    location_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS item_image (
    id BIGSERIAL PRIMARY KEY,
    item_id BIGINT NOT NULL,
    image_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS related_item (
    id BIGSERIAL PRIMARY KEY,
    item_id BIGINT NOT NULL,
    related_item_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS manufacturer (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    logo_url TEXT,
    description VARCHAR(255),
    group_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS location (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    image_url TEXT,
    description VARCHAR(255),
    group_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- this will allow to share entities with other users (add family members)
CREATE TABLE IF NOT EXISTS groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    group_owner_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_groups (
    user_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, group_id)
);

-- Add foreign key constraint to groups table referencing users table
ALTER TABLE groups
ADD CONSTRAINT fk_group_owner
FOREIGN KEY (group_owner_id)
REFERENCES users(id);

-- Add foreign key constraint to users table referencing groups table
ALTER TABLE users
ADD CONSTRAINT fk_active_group
FOREIGN KEY (active_group_id)
REFERENCES groups(id);

-- Add foreign key constraint to user_groups table referencing users table
ALTER TABLE user_groups 
ADD CONSTRAINT fk_user
FOREIGN KEY (user_id)
REFERENCES users(id);

-- Add foreign key constraint to user_groups table referencing groups table
ALTER TABLE user_groups 
ADD CONSTRAINT fk_group
FOREIGN KEY (group_id)
REFERENCES groups(id);

-- Add foreign key constraint to item table referencing users table
ALTER TABLE item
ADD CONSTRAINT fk_item_user_id
FOREIGN KEY (user_id)
REFERENCES users(id);

-- Add foreign key constraint to item table referencing group table
ALTER TABLE item
ADD CONSTRAINT fk_item_group_id
FOREIGN KEY (group_id)
REFERENCES groups(id);

-- Add foreign key constraint to location table referencing group table
ALTER TABLE location
ADD CONSTRAINT fk_location_group_id
FOREIGN KEY (group_id)
REFERENCES groups(id);

-- Add foreign key constraint to manufacturer table referencing group table
ALTER TABLE manufacturer
ADD CONSTRAINT fk_manufacturer_group_id
FOREIGN KEY (group_id)
REFERENCES groups(id);

-- Add foreign key constraint to item table referencing item_type table
ALTER TABLE item
ADD CONSTRAINT fk_item_item_type_id
FOREIGN KEY (item_type_id)
REFERENCES item_type(id);

-- Add foreign key constraint to item_info table referencing item table
ALTER TABLE item_info
ADD CONSTRAINT fk_item_info_item_id
FOREIGN KEY (item_id)
REFERENCES item(id);

-- Add foreign key constraint to item_image table referencing item table
ALTER TABLE item_image
ADD CONSTRAINT fk_item_image_item_id
FOREIGN KEY (item_id)
REFERENCES item(id);

-- Add foreign key constraint to related_item table referencing item table
ALTER TABLE related_item
ADD CONSTRAINT fk_related_item_item_id
FOREIGN KEY (item_id)
REFERENCES item(id);

-- Add foreign key constraint to related_item table referencing item table
ALTER TABLE related_item
ADD CONSTRAINT fk_related_item_related_item_id
FOREIGN KEY (related_item_id)
REFERENCES item(id);

-- Add foreign key constraint to item_info table referencing location table
ALTER TABLE item_info
ADD CONSTRAINT fk_item_info_location_id
FOREIGN KEY (location_id)
REFERENCES location(id);