-- init.sql
DROP TABLE IF EXISTS stock;
CREATE TABLE stock (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    cantidad INT NOT NULL,
    unidad VARCHAR(255) NOT NULL
);
INSERT INTO stock (nombre, cantidad, unidad) VALUES ('Harina', 10000, 'gr');
INSERT INTO stock (nombre, cantidad, unidad) VALUES ('Azúcar', 7500, 'gr');
INSERT INTO stock (nombre, cantidad, unidad) VALUES ('Levadura', 20, 'gr');
INSERT INTO stock (nombre, cantidad, unidad) VALUES ('Huevos', 48, 'u');
INSERT INTO stock (nombre, cantidad, unidad) VALUES ('Leche', 20000, 'ml');