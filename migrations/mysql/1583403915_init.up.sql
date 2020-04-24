CREATE TABLE `users` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `verification_token` char(64),
  `hashed_password` varchar(255) NOT NULL,
  `password_reset_token` char(64),
  `state` int(11) NOT NULL DEFAULT 2,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  UNIQUE KEY `idx_users_email` (`email`),
  UNIQUE KEY `idx_users_verification_token` (`verification_token`),
  UNIQUE KEY `idx_users_password_reset_token` (`password_reset_token`),
  KEY `idx_users_state` (`state`),
  KEY `idx_users_deleted_at` (`deleted_at`),
  KEY `idx_users_created_at` (`created_at`),
  KEY `idx_users_updated_at` (`updated_at`)
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

CREATE TABLE `posts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `state` INT NOT NULL DEFAULT 1,
  `user_id` BIGINT NOT NULL,
  `title` VARCHAR(64) NOT NULL,
  `content` TEXT NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_posts_user_id` (`user_id`),
  KEY `idx_posts_state` (`state`),
  KEY `idx_posts_created_at` (`created_at`),
  KEY `idx_posts_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO`users`(id, username, email, hashed_password, `state`, created_at) VALUES
(1, 'admin', 'admin@example.com', '$2a$12$R/Agn3zMt2iDF2/VBduy7uR1QLBoSeWrrCEgWByVFDsbRCl6Etbk2', 2, 0);