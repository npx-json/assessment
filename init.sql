CREATE TABLE IF NOT EXISTS whitelisted_tb (
    id INT AUTO_INCREMENT PRIMARY KEY,
    country VARCHAR(255)
);

INSERT INTO whitelisted_tb (id, country) VALUES
(1, 'US'),
(2, 'UK'),
(3, 'AU'),
(4, 'JP');