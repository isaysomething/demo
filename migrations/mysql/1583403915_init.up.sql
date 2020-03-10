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
) ENGINE=InnoDB;

CREATE TABLE `sessions` (
  `token` varchar(255) NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  `ip_address` varchar(255) NOT NULL DEFAULT '',
  `user_agent` varchar(255) NOT NULL DEFAULT '',
  `expired_at` datetime NOT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`token`),
  KEY `idx_sessions_user_id` (`user_id`),
  KEY `idx_sessions_expired_at` (`expired_at`)
) ENGINE=InnoDB;