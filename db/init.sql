-- init.sql
DROP TABLE stock
CREATE TABLE stock  (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    cantidad INT NOT NULL
);

INSERT INTO stock (nombre, cantidad) VALUES ('Harina', 100);
INSERT INTO stock (nombre, cantidad) VALUES ('Azúcar', 50);
INSERT INTO stock (nombre, cantidad) VALUES ('Levadura', 20);
