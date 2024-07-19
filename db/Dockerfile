# Use the official MySQL image
FROM mysql:latest

# Configuración de variables de entorno para MySQL
ENV MYSQL_ROOT_PASSWORD=rootagos
ENV MYSQL_DATABASE=panaderia_stock
# ENV MYSQL_USER=root
# ENV MYSQL_PASSWORD=rootagos

# Expose port 3306 to the outside world
EXPOSE 3306

# Copy the SQL script to initialize the database
COPY ./init.sql /docker-entrypoint-initdb.d/
