DROP TRIGGER IF EXISTS trg_update_latest_transaction ON transactions;
DROP FUNCTION IF EXISTS update_latest_transaction();

DROP TABLE IF EXISTS public.transactions_latest;
DROP TABLE IF EXISTS public.transactions;
DROP TABLE IF EXISTS public.account_categories;
DROP TABLE IF EXISTS public.merchants;
DROP TABLE IF EXISTS public.mccs;
DROP TABLE IF EXISTS public.accounts;
DROP TABLE IF EXISTS public.categories;