FROM postgis/postgis:16-3.4

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres

COPY ./ddl.sql /docker-entrypoint-initdb.d/001.sql
COPY ./dml.sql /docker-entrypoint-initdb.d/002.sql
COPY ./cities-dml.sql /docker-entrypoint-initdb.d/003.sql
COPY ./drugs-dml.sql /docker-entrypoint-initdb.d/004.sql
COPY ./pharmacy-dml.sql /docker-entrypoint-initdb.d/005.sql
COPY ./pharmacy-geom.sql /docker-entrypoint-initdb.d/006.sql
COPY ./pharmacy-drugs-dml-1.sql /docker-entrypoint-initdb.d/007.sql
# COPY ./pharmacy-drugs-dml-2.sql /docker-entrypoint-initdb.d/008.sql
# COPY ./pharmacy-drugs-dml-3.sql /docker-entrypoint-initdb.d/009.sql
# COPY ./pharmacy-drugs-dml-4.sql /docker-entrypoint-initdb.d/010.sql
# COPY ./pharmacy-drugs-dml-5.sql /docker-entrypoint-initdb.d/011.sql
# COPY ./pharmacy-drugs-dml-6.sql /docker-entrypoint-initdb.d/012.sql
# COPY ./pharmacy-drugs-dml-7.sql /docker-entrypoint-initdb.d/013.sql
# COPY ./pharmacy-drugs-dml-8.sql /docker-entrypoint-initdb.d/014.sql
# COPY ./pharmacy-drugs-dml-9.sql /docker-entrypoint-initdb.d/015.sql
# COPY ./pharmacy-drugs-dml-10.sql /docker-entrypoint-initdb.d/016.sql
# COPY ./pharmacy-drugs-dml-11.sql /docker-entrypoint-initdb.d/017.sql
# COPY ./pharmacy-drugs-dml-12.sql /docker-entrypoint-initdb.d/018.sql
# COPY ./pharmacy-drugs-dml-13.sql /docker-entrypoint-initdb.d/019.sql
# COPY ./pharmacy-drugs-dml-14.sql /docker-entrypoint-initdb.d/020.sql
# COPY ./pharmacy-drugs-dml-15.sql /docker-entrypoint-initdb.d/021.sql
# COPY ./pharmacy-drugs-dml-16.sql /docker-entrypoint-initdb.d/022.sql
# COPY ./pharmacy-drugs-dml-17.sql /docker-entrypoint-initdb.d/023.sql
# COPY ./pharmacy-drugs-dml-18.sql /docker-entrypoint-initdb.d/024.sql
# COPY ./pharmacy-drugs-dml-19.sql /docker-entrypoint-initdb.d/025.sql
# COPY ./pharmacy-drugs-dml-20.sql /docker-entrypoint-initdb.d/026.sql

HEALTHCHECK CMD pg_isready -U postgres || exit 1