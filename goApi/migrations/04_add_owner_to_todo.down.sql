alter table todos
    drop constraint if exists todos_users_id_fk;
alter table todos
    drop column if exists owner_id;