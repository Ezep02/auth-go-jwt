CREATE TABLE users (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(255) NOT NULL UNIQUE,
  "age" int NOT NULL,
  "email" varchar(255) NOT NULL UNIQUE,
  "password" varchar(255) NOT NULL,
  "admin" boolean NOT NULL,
  "is_user" boolean NOT NULL,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
  "id" varchar(255) PRIMARY KEY NOT NULL,
  "user_email" varchar(255) NOT NULL,
  "refresh_token" varchar(512) NOT NULL,
  "is_revoked" boolean NOT NULL DEFAULT false,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "expires_at" timestamp
);
