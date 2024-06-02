-- Remove foreign key constraint from user_groups table
ALTER TABLE user_groups DROP CONSTRAINT IF EXISTS fk_user;
ALTER TABLE user_groups DROP CONSTRAINT IF EXISTS fk_group;

-- Remove foreign key constraint from users table
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_active_group;

-- Remove foreign key constraint from groups table
ALTER TABLE groups DROP CONSTRAINT IF EXISTS fk_group_owner;

-- Remove foreign key constraint from location table
ALTER TABLE location DROP CONSTRAINT IF EXISTS fk_location_group_id;

-- Remove foreign key constraint from manufacturer table
ALTER TABLE manufacturer DROP CONSTRAINT IF EXISTS fk_manufacturer_group_id;

-- Remove foreign key constraint from item_info table
ALTER TABLE item_info DROP CONSTRAINT IF EXISTS fk_item_info_location_id;

-- Remove foreign key constraint from related_item table
ALTER TABLE related_item DROP CONSTRAINT IF EXISTS fk_related_item_related_item_id;
ALTER TABLE related_item DROP CONSTRAINT IF EXISTS fk_related_item_item_id;

-- Remove foreign key constraint from item_image table
ALTER TABLE item_image DROP CONSTRAINT IF EXISTS fk_item_image_item_id;

-- Remove foreign key constraint from item_info table
ALTER TABLE item_info DROP CONSTRAINT IF EXISTS fk_item_info_item_id;

-- Remove foreign key constraint from item table
ALTER TABLE item DROP CONSTRAINT IF EXISTS fk_item_item_type_id;
ALTER TABLE item DROP CONSTRAINT IF EXISTS fk_item_user_id;
ALTER TABLE item DROP CONSTRAINT IF EXISTS fk_item_group_id;

-- Drop tables

-- Drop invitations table
DROP TABLE IF EXISTS invitations;

-- Drop groups table
DROP TABLE IF EXISTS groups;

-- Drop user_groups table
DROP TABLE IF EXISTS user_groups;

-- Drop manufacturer table
DROP TABLE IF EXISTS manufacturer;

-- Drop location table
DROP TABLE IF EXISTS location;

-- Drop related_item table
DROP TABLE IF EXISTS related_item;

-- Drop item_image table
DROP TABLE IF EXISTS item_image;

-- Drop item_info table
DROP TABLE IF EXISTS item_info;

-- Drop item_type table
DROP TABLE IF EXISTS item_type;

-- Drop item table
DROP TABLE IF EXISTS item;

-- Drop users table
DROP TABLE IF EXISTS users;