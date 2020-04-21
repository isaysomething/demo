INSERT INTO `auth_rules`(`id`, `p_type`, `v0`, `v1`, `v2`) VALUES
(null, 'p', 'user', '/v1/user/info', 'GET'),
(null, 'p', 'admin', '/v1/user/info', 'GET'),

(null, 'p', 'admin', '/v1/posts', 'GET|POST'),
(null, 'p', 'admin', '/v1/posts/*', 'GET|PUT|DELETE'),

(null, 'p', 'admin', '/v1/users', 'GET|POST'),
(null, 'p', 'admin', '/v1/users/*', 'GET|PUT|DELETE'),

(null, 'g', 'user_1', 'admin', ''),
(null, 'g', 'user_1', 'user', '');