CREATE TABLE IF NOT EXISTS public.devices
(
    id              serial PRIMARY KEY,
    organization_id         integer NOT NULL references public.organizations(id),
    room_id  integer references public.rooms(id),
    "guid"         UUID,
    inventory_number VARCHAR(255),
    serial_number   VARCHAR(255),
    characteristics TEXT,
    category        VARCHAR(255),
    units           VARCHAR(255),
    power_consumption FLOAT,
    created_date    timestamptz NOT NULL,
    updated_date    timestamptz NOT NULL,
    deleted_date    timestamptz
);