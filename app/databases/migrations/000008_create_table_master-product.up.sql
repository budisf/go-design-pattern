CREATE TABLE IF NOT EXISTS public.products (
    product_id bigserial NOT NULL,
    product_code int8 NULL,
    "name" text NULL,
    price_value int8 NULL,
    incentive_percentage numeric NULL,
    period_sales timestamptz NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    CONSTRAINT products_pkey PRIMARY KEY (product_id)
);