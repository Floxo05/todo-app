alter table todos
    add owner_id int not null,
    add constraint todos_users_id_fk
        foreign key (owner_id) references users (id);

