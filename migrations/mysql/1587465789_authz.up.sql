INSERT INTO `auth_rules`(`id`, `p_type`, `v0`, `v1`, `v2`) VALUES
(null, 'p', 'user', '/v1/user/info', 'GET'),
(null, 'p', 'admin', '/v1/user/info', 'GET'),

(null, 'p', 'admin', '/v1/posts', 'GET|POST'),
(null, 'p', 'admin', '/v1/posts/*', 'GET|PUT|PATCH|DELETE'),

(null, 'p', 'admin', '/v1/users', 'GET|POST'),
(null, 'p', 'admin', '/v1/users/*', 'GET|PUT|PATCH|DELETE'),

(null, 'p', 'admin', '/v1/roles', 'GET|POST'),
(null, 'p', 'admin', '/v1/roles/*', 'GET|PUT|PATCH|DELETE'),

(null, 'g', 'user_1', 'admin', ''),
(null, 'g', 'user_1', 'user', '');

INSERT INTO `auth_item_groups`(`id`, `name`, `created_at`, `updated_at`) VALUES
(1, 'User management', NOW(), NOW()),
(2, 'Post management', NOW(), NOW()),
(3, 'Role management', NOW(), NOW()),
(4, 'Others', NOW(), NOW());

INSERT INTO `auth_items`(`id`, `name`, `item_type`, `reserved`, `group_id`, `obj`, `act`, `created_at`, `updated_at`) VALUES
('admin', 'Administrator', 1, 1, 0, '', '', NOW(), NOW()),
('user', 'User', 1, 1, 0, '', '', NOW(), NOW()),

('user:info', 'User info', 2, 1, 3, '/v1/user/info', 'GET', NOW(), NOW()),

('post:list', 'List posts', 2, 1, 2, '/v1/posts', 'GET', NOW(), NOW()),
('post:create', 'Create post', 2, 1, 2, '/v1/posts', 'POST', NOW(), NOW()),
('post:view', 'View Post', 2, 1, 2, '/v1/posts/*', 'GET', NOW(), NOW()),
('post:edit', 'Edit post', 2, 1, 2, '/v1/posts/*', 'PUT|PATCH', NOW(), NOW()),
('post:delete', 'Delete post', 2, 1, 2, '/v1/posts/*', 'DELETE', NOW(), NOW()),

('user:list', 'List users', 2, 1, 1, '/v1/users', 'GET', NOW(), NOW()),
('user:create', 'Create user', 2, 1, 1, '/v1/users', 'POST', NOW(), NOW()),
('user:view', 'View user', 2, 1, 1, '/v1/users/*', 'GET', NOW(), NOW()),
('user:edit', 'Edit user', 2, 1, 1, '/v1/users/*', 'PUT|PATCH', NOW(), NOW()),
('user:delete', 'Delete user', 2, 1, 1, '/v1/users/*', 'DELETE', NOW(), NOW()),

('role:list', 'List roles', 2, 1, 3, '/v1/roles', 'GET', NOW(), NOW()),
('role:create', 'Create role', 2, 1, 3, '/v1/roles', 'POST', NOW(), NOW()),
('role:view', 'View role', 2, 1, 3, '/v1/roles/*', 'GET', NOW(), NOW()),
('role:edit', 'Edit role', 2, 1, 3, '/v1/roles/*', 'PUT|PATCH', NOW(), NOW()),
('role:delete', 'Delete role', 2, 1, 3, '/v1/roles/*', 'DELETE', NOW(), NOW());