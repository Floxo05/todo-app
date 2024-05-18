CREATE TABLE user_todos
(
    user_id INT,
    todo_id INT,
    PRIMARY KEY (user_id, todo_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (todo_id) REFERENCES todos (id)
);