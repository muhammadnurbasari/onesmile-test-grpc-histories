CREATE TABLE histories (
  "id" serial NOT NULL,
  "credit_card" varchar(50) COLLATE "pg_catalog"."default",
  "grand_total" int,
  CONSTRAINT "histories_pkey" PRIMARY KEY ("id")
);