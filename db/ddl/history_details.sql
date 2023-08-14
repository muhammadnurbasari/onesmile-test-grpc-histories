CREATE TABLE history_details (
  "id" serial NOT NULL,
  "history_id" int,
  "name" varchar(50) COLLATE "pg_catalog"."default",
  "quantity" int,
  "sub_total" int,
  CONSTRAINT "history_details_pkey" PRIMARY KEY ("id")
);