ALTER TABLE todos
    ADD category_id  INT,
    ADD CONSTRAINT todos_categories_id_fk
        FOREIGN KEY (category_id) REFERENCES categories (id);