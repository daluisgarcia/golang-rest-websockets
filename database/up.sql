DROP TABLE IF EXISTS "users";

CREATE TABLE "users" (
  id varchar(36) NOT NULL PRIMARY KEY,
  email varchar(255) UNIQUE NOT NULL,
  password varchar(255) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS "posts";

CREATE TABLE posts (
 	id varchar(36) NOT NULL PRIMARY KEY,
	post_content varchar(32) NOT NULL,
	user_id varchar(36) NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id)
)
