-- Drop foreign key constraint from parent_id column
ALTER TABLE item_type DROP CONSTRAINT fk_parent_item_type;

-- Remove parent_id column from item_type table
ALTER TABLE item_type DROP COLUMN parent_id;
