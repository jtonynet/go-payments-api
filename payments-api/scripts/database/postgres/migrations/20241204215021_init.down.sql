DROP TABLE IF EXISTS public.transactions CASCADE;
DROP TABLE IF EXISTS public.account_categories CASCADE;
DROP TABLE IF EXISTS public.merchants CASCADE;
DROP TABLE IF EXISTS public.mccs CASCADE;
DROP TABLE IF EXISTS public.accounts CASCADE;
DROP TABLE IF EXISTS public.categories CASCADE;

ALTER SEQUENCE public.transactions_id_seq RESTART WITH 1;
ALTER SEQUENCE public.account_categories_id_seq RESTART WITH 1;
ALTER SEQUENCE public.merchants_id_seq RESTART WITH 1;
ALTER SEQUENCE public.mccs_id_seq RESTART WITH 1;
ALTER SEQUENCE public.accounts_id_seq RESTART WITH 1;
ALTER SEQUENCE public.categories_id_seq RESTART WITH 1;
