CREATE TABLE public.accounts (
    id bigserial NOT NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    uid uuid NULL,
    "name" varchar(255) NULL,
    CONSTRAINT accounts_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_accounts_deleted_at ON public.accounts USING btree (deleted_at);
CREATE UNIQUE INDEX idx_accounts_uid ON public.accounts USING btree (uid);

CREATE TABLE public.categories (
    id bigserial NOT NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    uid uuid NULL,
    "name" varchar(255) NULL,
    priority int8 NULL,
    CONSTRAINT categories_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_categories_deleted_at ON public.categories USING btree (deleted_at);
CREATE UNIQUE INDEX idx_categories_uid ON public.categories USING btree (uid);

CREATE TABLE public.account_categories (
    id bigserial NOT NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    account_id int8 NULL,
    category_id int8 NULL,
    CONSTRAINT account_categories_pkey PRIMARY KEY (id),
    CONSTRAINT fk_accounts_account_categories FOREIGN KEY (account_id) REFERENCES public.accounts(id),
    CONSTRAINT fk_categories_account_categories FOREIGN KEY (category_id) REFERENCES public.categories(id)
);
CREATE INDEX idx_account_categories_deleted_at ON public.account_categories USING btree (deleted_at);

CREATE TABLE public.mccs (
    id bigserial NOT NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    uid uuid NULL,
    category_id int8 NULL,
    mcc varchar(5) NULL,
    CONSTRAINT mccs_pkey PRIMARY KEY (id),
    CONSTRAINT fk_categories_mc_cs FOREIGN KEY (category_id) REFERENCES public.categories(id)
);
CREATE INDEX idx_mccs_deleted_at ON public.mccs USING btree (deleted_at);
CREATE UNIQUE INDEX idx_mccs_uid ON public.mccs USING btree (uid);

CREATE TABLE public.merchants (
    id bigserial NOT NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    uid uuid NULL,
    "name" varchar(255) NULL,
    mcc_id int8 NULL,
    CONSTRAINT merchants_pkey PRIMARY KEY (id),
    CONSTRAINT fk_mccs_merchants FOREIGN KEY (mcc_id) REFERENCES public.mccs(id)
);
CREATE INDEX idx_merchants_deleted_at ON public.merchants USING btree (deleted_at);
CREATE UNIQUE INDEX idx_merchants_name ON public.merchants USING btree ("name");
CREATE UNIQUE INDEX idx_merchants_uid ON public.merchants USING btree (uid);

CREATE TABLE public.transactions (
    id bigserial NOT NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    uid uuid NULL,
    account_id int8 NULL,
    category_id int8 NULL,
    amount numeric(20, 2) NULL,
    mcc varchar(5) NULL,
    merchant_name varchar(255) NULL,
    CONSTRAINT transactions_pkey PRIMARY KEY (id),
    CONSTRAINT uni_transactions_uid UNIQUE (uid),
    CONSTRAINT fk_categories_transactions FOREIGN KEY (category_id) REFERENCES public.categories(id),
    CONSTRAINT fk_transactions_account FOREIGN KEY (account_id) REFERENCES public.accounts(id)
);
CREATE INDEX idx_transaction_composite ON public.transactions USING btree (account_id, category_id);
CREATE INDEX idx_transactions_deleted_at ON public.transactions USING btree (deleted_at);
