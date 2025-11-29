CREATE TABLE "main"."users" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_name" character varying NOT NULL,
    "name" character varying NOT NULL,
    "password" character varying NOT NULL,
    "status" "main"."status_enum" NOT NULL DEFAULT 'NEW',
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    CONSTRAINT "pk_user" PRIMARY KEY ("id"),
    CONSTRAINT "uq_user_name" UNIQUE ("user_name")
);
