CREATE TABLE "main"."projects" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" character varying NOT NULL,
    "description" character varying NOT NULL,
    "status" "main"."status_enum" NOT NULL DEFAULT 'NEW',
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    CONSTRAINT "pk_project" PRIMARY KEY ("id"),
    CONSTRAINT "uq_project_name" UNIQUE ("name")
);
