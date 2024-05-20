CREATE TABLE IF NOT EXISTS public.rooms
(
    id              serial PRIMARY KEY,
    organization_id         integer NOT NULL references public.organizations(id),
    "name"          varchar(200) NOT NULL,
    "description"   text,
    created_date    timestamptz NOT NULL,
    updated_date    timestamptz NOT NULL,
    deleted_date    timestamptz
);