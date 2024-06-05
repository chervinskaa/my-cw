CREATE TABLE IF NOT EXISTS public.events
(
    id              serial PRIMARY KEY,
    device_id         integer NOT NULL references public.devices(id),
    room_id  integer references public.rooms(id),
    "action" VARCHAR(50) NOT NULL,
    created_date    timestamptz NOT NULL,
    updated_date    timestamptz NOT NULL,
    deleted_date    timestamptz
);