# Build stage
FROM golang:1.22.0-alpine AS builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy and download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Copy the GeoLite2-Country.mmdb into the container
COPY GeoLite2-Country.mmdb ./GeoLite2-Country.mmdb

# Build the Go application
RUN go build -o ipcheck ./main.go

# Final stage
FROM alpine:3.17
WORKDIR /root/

# Install MySQL server and client
RUN apk add --no-cache mysql mysql-client openrc

# Copy the built application from the builder stage
COPY --from=builder /app/ipcheck .

# Copy the GeoLite2-Country.mmdb into the container
COPY GeoLite2-Country.mmdb ./GeoLite2-Country.mmdb

# Copy a MySQL configuration file (optional)
COPY my.cnf /etc/mysql/my.cnf

# Copy a MySQL initialization script
COPY mysql-init.sql /docker-entrypoint-initdb.d/mysql-init.sql

# Initialize MySQL data directory
RUN mkdir -p /var/lib/mysql /var/run/mysqld && \
    chown -R mysql:mysql /var/lib/mysql /var/run/mysqld && \
    mysql_install_db --user=mysql --datadir=/var/lib/mysql

# Expose ports for HTTP, gRPC, and MySQL
EXPOSE 8080 9090 3306



# Start both MySQL and the Go application
CMD ["sh", "-c", "mysqld_safe --datadir=/var/lib/mysql & ./ipcheck"]