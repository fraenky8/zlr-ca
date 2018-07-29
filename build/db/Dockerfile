FROM postgres

#
# Schema
#
COPY database.sql /docker-entrypoint-initdb.d/10-database.sql

#
# Data
#
COPY data/20-icecream.sql /docker-entrypoint-initdb.d/20-icecream.sql
COPY data/30-ingredients.sql /docker-entrypoint-initdb.d/30-ingredients.sql
COPY data/40-sourcing_values.sql /docker-entrypoint-initdb.d/40-sourcing_values.sql
COPY data/50-icecream_has_ingredients.sql /docker-entrypoint-initdb.d/50-icecream_has_ingredients.sql
COPY data/60-icecream_has_sourcing_values.sql /docker-entrypoint-initdb.d/60-icecream_has_sourcing_values.sql
