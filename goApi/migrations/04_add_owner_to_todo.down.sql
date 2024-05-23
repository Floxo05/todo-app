alter table todos
    drop constraint if exists todos_users_id_fk,
    drop column if exists owner_id;