CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `verification_token` varchar(255),
  `hashed_password` varchar(255) NOT NULL,
  `password_reset_token` varchar(255),
  `status` int(11) NOT NULL DEFAULT 10,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `verification_token` (`verification_token`),
  UNIQUE KEY `password_reset_token` (`password_reset_token`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `auth_rules` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `p_type` VARCHAR(32) NOT NULL DEFAULT '',
    `v0` VARCHAR(255) NOT NULL DEFAULT '',
    `v1` VARCHAR(255) NOT NULL DEFAULT '',
    `v2` VARCHAR(255) NOT NULL DEFAULT '',
    `v3` VARCHAR(255) NOT NULL DEFAULT '',
    `v4` VARCHAR(255) NOT NULL DEFAULT '',
    `v5` VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
