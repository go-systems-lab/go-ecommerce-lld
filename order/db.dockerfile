FROM postgres:17.3-alpine

# Set default database name
ENV POSTGRES_DB=ecommerce_order
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres

EXPOSE 5432

CMD ["postgres"]