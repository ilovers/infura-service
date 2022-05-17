CREATE TABLE `blocks` (
                          `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                          `created_at` datetime(3) DEFAULT NULL,
                          `updated_at` datetime(3) DEFAULT NULL,
                          `deleted_at` datetime(3) DEFAULT NULL,
                          `number` bigint(20) NOT NULL,
                          `hash` varchar(66) NOT NULL,
                          `parent_hash` varchar(66) DEFAULT NULL,
                          `transactions_root` varchar(66) DEFAULT NULL,
                          `state_root` varchar(66) DEFAULT NULL,
                          `miner` varchar(42) DEFAULT NULL,
                          `size` int(11) DEFAULT NULL,
                          `gas_limit` bigint(20) unsigned DEFAULT NULL,
                          `gas_used` bigint(20) unsigned DEFAULT NULL,
                          `timestamp` int(11) DEFAULT NULL,
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `unique_hash` (`hash`),
                          KEY `idx_blocks_deleted_at` (`deleted_at`),
                          KEY `idx_blocks_number` (`number`)
) ENGINE=InnoDB AUTO_INCREMENT=3534 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `contract_codes` (
                                  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                                  `created_at` datetime(3) DEFAULT NULL,
                                  `updated_at` datetime(3) DEFAULT NULL,
                                  `deleted_at` datetime(3) DEFAULT NULL,
                                  `address` varchar(42) NOT NULL,
                                  `code` longtext,
                                  `block_number` bigint(20) DEFAULT NULL,
                                  PRIMARY KEY (`id`),
                                  UNIQUE KEY `unique_address` (`address`),
                                  KEY `idx_contract_codes_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `log_topics` (
                              `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                              `created_at` datetime(3) DEFAULT NULL,
                              `updated_at` datetime(3) DEFAULT NULL,
                              `deleted_at` datetime(3) DEFAULT NULL,
                              `topic` varchar(66) DEFAULT NULL,
                              `transaction_log_id` bigint(20) unsigned DEFAULT NULL,
                              PRIMARY KEY (`id`),
                              KEY `idx_log_topics_deleted_at` (`deleted_at`),
                              KEY `fk_transaction_logs_topics` (`transaction_log_id`),
                              CONSTRAINT `fk_transaction_logs_topics` FOREIGN KEY (`transaction_log_id`) REFERENCES `transaction_logs` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `transaction_logs` (
                                    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                                    `created_at` datetime(3) DEFAULT NULL,
                                    `updated_at` datetime(3) DEFAULT NULL,
                                    `deleted_at` datetime(3) DEFAULT NULL,
                                    `address` varchar(42) NOT NULL,
                                    `data` text,
                                    `transaction_hash` varchar(66) DEFAULT NULL,
                                    `transaction_index` int(11) DEFAULT NULL,
                                    `log_index` int(11) DEFAULT NULL,
                                    `block_hash` varchar(66) NOT NULL,
                                    `block_number` bigint(20) NOT NULL,
                                    `transaction_receipt_id` bigint(20) unsigned DEFAULT NULL,
                                    PRIMARY KEY (`id`),
                                    KEY `idx_transaction_logs_deleted_at` (`deleted_at`),
                                    KEY `idx_transaction_logs_address` (`address`),
                                    KEY `idx_transaction_logs_block_hash` (`block_hash`),
                                    KEY `idx_transaction_logs_block_number` (`block_number`),
                                    KEY `fk_transaction_receipts_logs` (`transaction_receipt_id`),
                                    CONSTRAINT `fk_transaction_receipts_logs` FOREIGN KEY (`transaction_receipt_id`) REFERENCES `transaction_receipts` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `transaction_receipts` (
                                        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                                        `created_at` datetime(3) DEFAULT NULL,
                                        `updated_at` datetime(3) DEFAULT NULL,
                                        `deleted_at` datetime(3) DEFAULT NULL,
                                        `status` tinyint(4) DEFAULT NULL,
                                        `cumulative_gas_used` int(11) DEFAULT NULL,
                                        `transaction_hash` varchar(66) NOT NULL,
                                        `contract_address` varchar(42) DEFAULT NULL,
                                        `gas_used` int(11) DEFAULT NULL,
                                        `block_hash` varchar(66) DEFAULT NULL,
                                        `block_number` bigint(20) DEFAULT NULL,
                                        `transaction_index` int(11) DEFAULT NULL,
                                        `from` varchar(42) DEFAULT NULL,
                                        `to` varchar(42) DEFAULT NULL,
                                        PRIMARY KEY (`id`),
                                        UNIQUE KEY `unique_hash` (`transaction_hash`),
                                        KEY `idx_transaction_receipts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `transactions` (
                                `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                                `created_at` datetime(3) DEFAULT NULL,
                                `updated_at` datetime(3) DEFAULT NULL,
                                `deleted_at` datetime(3) DEFAULT NULL,
                                `block_hash` varchar(66) DEFAULT NULL,
                                `block_number` bigint(20) DEFAULT NULL,
                                `from` varchar(42) DEFAULT NULL,
                                `gas` int(11) DEFAULT NULL,
                                `gas_price` varchar(66) DEFAULT NULL,
                                `hash` varchar(66) DEFAULT NULL,
                                `input` text,
                                `nonce` int(11) DEFAULT NULL,
                                `to` varchar(42) DEFAULT NULL,
                                `index` int(11) DEFAULT NULL,
                                `value` varchar(255) DEFAULT NULL,
                                `v` varchar(255) DEFAULT NULL,
                                `r` varchar(255) DEFAULT NULL,
                                `s` varchar(255) DEFAULT NULL,
                                `block_id` bigint(20) unsigned DEFAULT NULL,
                                PRIMARY KEY (`id`),
                                KEY `idx_transactions_deleted_at` (`deleted_at`),
                                KEY `fk_blocks_transactions` (`block_id`),
                                CONSTRAINT `fk_blocks_transactions` FOREIGN KEY (`block_id`) REFERENCES `blocks` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;