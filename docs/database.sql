--
-- Schema zlr_ca
--
DROP SCHEMA IF EXISTS zlr_ca CASCADE;
CREATE SCHEMA IF NOT EXISTS zlr_ca;

--
-- Table icecream
--
create table if not exists zlr_ca.icecream
(
  product_id             integer      not null
    constraint icecream_pkey
    primary key,
  name                   varchar(200) not null,
  description            varchar(200),
  story                  text,
  image_open             varchar(200),
  image_closed           varchar(200),
  allergy_info           varchar(200),
  dietary_certifications varchar(50)
);
create unique index if not exists icecream_product_id_uindex
  on zlr_ca.icecream (product_id);

--
-- Table ingredients
--
create table zlr_ca.ingredients
(
  id   serial      not null,
  name varchar(50) not null,
  constraint ingredients_id_name_pk
  primary key (id, name)
);

create unique index ingredients_id_uindex
  on zlr_ca.ingredients (id);

create unique index ingredients_name_uindex
  on zlr_ca.ingredients (name);

--
-- Table sourcing_values
--
create table zlr_ca.sourcing_values
(
  id          serial       not null,
  description varchar(200) not null,
  constraint sourcing_values_id_description_pk
  primary key (id, description)
);

create unique index sourcing_values_id_uindex
  on zlr_ca.sourcing_values (id);

create unique index sourcing_values_description_uindex
  on zlr_ca.sourcing_values (description);

--
-- Table icecream_has_ingredients
--
create table zlr_ca.icecream_has_ingredients
(
  icecream_product_id integer not null
    constraint icecream_has_ingredients_icecream_product_id_fk
    references zlr_ca.icecream (product_id)
    on delete cascade,
  ingredients_id      integer not null
    constraint icecream_has_ingredients_ingredients_id_fk
    references zlr_ca.ingredients (id)
    on delete cascade,
  constraint icecream_has_ingredients_icecream_product_id_ingredients_id_pk
  primary key (icecream_product_id, ingredients_id)
);

--
-- Table icecream_has_sourcing_values
--
create table zlr_ca.icecream_has_sourcing_values
(
  icecream_product_id integer not null
    constraint icecream_has_sourcing_values_icecream_product_id_fk
    references zlr_ca.icecream (product_id)
    on delete cascade,
  sourcing_values_id  integer not null
    constraint icecream_has_sourcing_values_sourcing_values_id_fk
    references zlr_ca.sourcing_values (id)
    on delete cascade,
  constraint icecream_has_sourcing_values_icecream_id_sourcing_values_id_pk
  primary key (icecream_product_id, sourcing_values_id)
);

