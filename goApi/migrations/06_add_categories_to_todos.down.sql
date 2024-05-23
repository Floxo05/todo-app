ALTER TABLE todos
    DROP CONSTRAINT todos_categories_id_fk,
    DROP COLUMN IF EXISTS category_id;