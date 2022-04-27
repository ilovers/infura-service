CREATE TABLE `transaction_receipts` (
                                        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                                        `status` tinyint(4) unsigned DEFAULT NULL,
                                        `cumulative_gas_used` int(11) unsigned DEFAULT NULL,
                                        `transaction_hash` varchar(66) NOT NULL,
                                        `contract_address` varchar(42) DEFAULT NULL,
                                        `gas_used` int(11) unsigned DEFAULT NULL,
                                        `block_hash` varchar(66) DEFAULT NULL,
                                        `block_number` bigint(20) DEFAULT NULL,
                                        `transaction_index` int(11) unsigned DEFAULT NULL,
                                        `from` varchar(42) DEFAULT NULL,
                                        `to` varchar(42) DEFAULT NULL,
                                        PRIMARY KEY (`id`),
                                        UNIQUE KEY `unique_hash` (`transaction_hash`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `transaction_logs` (
                                    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                                    `address` varchar(42) NOT NULL,
                                    `data` varchar(256) DEFAULT NULL,
                                    `transaction_hash` varchar(66) NOT NULL,
                                    `transaction_index` int(11) unsigned DEFAULT NULL,
                                    `log_index` int(11) unsigned NOT NULL,
                                    `block_hash` varchar(66) NOT NULL,
                                    `block_number` bigint(20) NOT NULL,
                                    PRIMARY KEY (`id`),
                                    KEY `idx_transaction_logs_address` (`address`),
                                    KEY `idx_transaction_logs_transaction_hash` (`transaction_hash`),
                                    KEY `idx_transaction_logs_block_hash` (`block_hash`),
                                    KEY `idx_transaction_logs_block_number` (`block_number`),
                                    KEY `log_index` (`log_index`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `log_topics` (
                              `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                              `transaction_hash` varchar(66) NOT NULL,
                              `log_index` int(11) unsigned NOT NULL,
                              `topic` varchar(66) NOT NULL,
                              PRIMARY KEY (`id`),
                              KEY `idx_log_topics_transaction_hash` (`transaction_hash`),
                              KEY `idx_log_topics_topic` (`topic`),
                              KEY `log_index` (`log_index`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;