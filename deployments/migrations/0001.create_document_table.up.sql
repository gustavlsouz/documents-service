create table
  document (
    "id" uuid primary key default gen_random_uuid (),
    "type" varchar(4) not null,
    "value" varchar(14) not null,
    "createdAt" timestamp without time zone not null default now (),
    "updatedAt" timestamp without time zone not null default now (),
    "isBlocked" boolean default false,
    unique ("type", "value")
  );

CREATE INDEX "document_value_idx" ON "document" ("value");

CREATE INDEX "document_type_idx" ON "document" ("type");

CREATE INDEX "document_created_at_idx" ON "document" ("createdAt");

CREATE INDEX "document_updated_at_idx" ON "document" ("updatedAt");