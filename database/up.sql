DROP TABLE IF EXISTS "users";

CREATE TABLE "users" (
  id varchar(36) NOT NULL PRIMARY KEY,
  email varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
