-- this migration attempts to support nested item_type. 
-- in this way users will be able to categorise their items
-- more structuredly

-- Add parent_id column to item_type table
ALTER TABLE item_type ADD COLUMN parent_id BIGINT;

-- Add foreign key constraint to parent_id column
ALTER TABLE item_type
ADD CONSTRAINT fk_parent_item_type
FOREIGN KEY (parent_id) REFERENCES item_type (id);
