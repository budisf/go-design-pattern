CREATE TABLE IF NOT EXISTS public.products_bridging (
    id int8 NOT NULL,
    product_id int8 NULL,
    product_id_distributor int8 NOT NULL,
    product_name_distributor varchar(256) NOT NULL,
    distributor_id int8 NULL,
    created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    CONSTRAINT products_bridging_pkey PRIMARY KEY (id),
    CONSTRAINT fk_distributor FOREIGN KEY (distributor_id) REFERENCES public.distributors(id) ON UPDATE CASCADE ON DELETE NO ACTION,
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES public.products(product_id) ON UPDATE CASCADE ON DELETE NO ACTION
);