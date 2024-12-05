INSERT INTO categories (uid, name, priority, created_at, updated_at)
VALUES
    ('5681b4b5-6176-498a-a856-8932f79c05cc', 'FOOD', 1, NOW(), NOW()),
    ('7bcfcd2a-2fde-4564-916b-92410e794272', 'MEAL', 2, NOW(), NOW()),
    ('056de185-bff0-4c4a-93fa-7245f9e72b67', 'CASH', 3, NOW(), NOW());

INSERT INTO mccs (uid, mcc, category_id, created_at, updated_at)
VALUES
    ('11f0c06e-0dff-4643-86bf-998d11e9374f', '5411', 1, NOW(), NOW()),
    ('fe5a4c17-a7cd-4072-a793-e99e2642e21a', '5412', 1, NOW(), NOW()),
    ('5268ec2b-aa14-4d55-906a-13c91d89826c', '5811', 2, NOW(), NOW()),
    ('6179e57c-e630-4e2f-a5db-d153e0cdb9a9', '5812', 2, NOW(), NOW());

INSERT INTO merchants (uid, name, mcc_id, created_at, updated_at)
VALUES
    ('95abe1ff-6f67-4a17-a4eb-d4842e324f1f', 'UBER EATS                   SAO PAULO BR', 2, NOW(), NOW()),
    ('a53c6a52-8a18-4e7d-8827-7f612233c7ec', 'PAG*JoseDaSilva          RIO DE JANEI BR', 4, NOW(), NOW());

INSERT INTO accounts (uid, name, created_at, updated_at) 
VALUES('123e4567-e89b-12d3-a456-426614174000', 'Jonh Doe', NOW(), NOW());

INSERT INTO account_categories (account_id, category_id, created_at, updated_at)
VALUES
    (1, 1, NOW(), NOW()),
    (1, 2, NOW(), NOW()),
    (1, 3, NOW(), NOW());

INSERT INTO transactions (account_id, amount, category_id, created_at, updated_at)
VALUES
    (1, 205.11, 1, NOW(), NOW()),
    (1, 110.22, 2, NOW(), NOW()),
    (1, 115.33, 3, NOW(), NOW());
