CREATE TABLE categories
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    created_user_id INT NOT NULL,
    FOREIGN KEY (created_user_id) REFERENCES users (id)
);