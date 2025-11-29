CREATE TABLE "main"."releases" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "tag" character varying NOT NULL,
    "app" character varying NOT NULL,
    "status" "main"."status_enum" NOT NULL DEFAULT 'NEW',
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    CONSTRAINT "pk_version" PRIMARY KEY ("id"),
    CONSTRAINT "uq_version_tag_app" UNIQUE ("tag", "app"),
    CONSTRAINT "fk_app" FOREIGN KEY ("app") REFERENCES "main"."apps" ("name")
);
