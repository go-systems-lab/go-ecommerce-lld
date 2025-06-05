FROM postgres:17.3-alpine

# Copy migration files to PostgreSQL initialization directory
COPY ./migrations/*.sql /docker-entrypoint-initdb.d/

# Set default database name
ENV POSTGRES_DB=ecommerce_account
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres

EXPOSE 5432

CMD ["postgres"]