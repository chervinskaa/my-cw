CREATE TABLE IF NOT EXISTS public.organizations
(
    id              serial PRIMARY KEY,
    user_id         integer NOT NULL references public.users(id),
    "name"          varchar(200) NOT NULL,
    "description"   text,
    city            varchar(100) NOT NULL,
    "address"       text NOT NULL,
    lat             double precision NOT NULL,
    lon             double precision NOT NULL,
    created_date    timestamptz NOT NULL,
    updated_date    timestamptz NOT NULL,
    deleted_date    timestamptz
);