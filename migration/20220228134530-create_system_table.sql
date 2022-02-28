
-- +migrate Up
CREATE TABLE `system` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '系統名稱',
    `system_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '系統類型',
    `tag` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '系統標籤',
    `email` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `address` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '地址',
    `tel` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '電話',
    `uuid` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
    `quota` int(11) NOT NULL COMMENT '次數',
    `ip_address` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP Address',
    `mac_address` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'MAC Address',
    `is_disable` tinyint(4) NOT NULL COMMENT '0:啟用 1:禁用',
    `principal` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '公司負責人',
    `salesman` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '業務',
    `salesman_phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '業務電話',
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_system_tag` (`tag`),
    UNIQUE KEY `uni_system_uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_account` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `system_id` int(11) NOT NULL COMMENT 'ref:system.id',
    `account` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '帳號',
    `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `email` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名字',
    `is_disable` tinyint(4) NOT NULL COMMENT '0:啟用 1:禁用',
    `verify_at` datetime DEFAULT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_system_id_account` (`system_id`,`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_account_role` (
    `sys_account_id` int(11) NOT NULL COMMENT 'ref:sys_account.id',
    `sys_role_id` int(11) NOT NULL COMMENT 'ref:sys_role.id',
    `created_at` datetime NOT NULL,
    UNIQUE KEY `uni_samr_sami_sri` (`sys_account_id`,`sys_role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_permission` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `system_id` int(11) NOT NULL COMMENT 'ref:system.id',
    `allow_api_path` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '允許的api url',
    `action` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '動作',
    `slug` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '標籤',
    `description` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '描述',
    `is_disable` tinyint(4) NOT NULL COMMENT '0:啟用 1:禁用',
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_sp_sid_slug` (`system_id`,`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_purchase` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `system_id` int(11) NOT NULL COMMENT 'ref:system.id',
    `tracking_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '追蹤單號',
    `quota` int(11) NOT NULL COMMENT '次數',
    `salesman` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IMQ 業務',
    `applied_at` datetime DEFAULT NULL COMMENT '申請日期',
    `created_at` datetime NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_role` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `sort` int(11) NOT NULL COMMENT '排序',
    `system_id` int(11) NOT NULL COMMENT 'ref:system.id',
    `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名稱',
    `display_name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '顯示名稱',
    `is_disable` tinyint(4) NOT NULL COMMENT '0:啟用 1:禁用',
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_sr_si_name` (`system_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_role_permission` (
    `sys_role_id` int(11) NOT NULL COMMENT 'ref:sys_role.id',
    `sys_permission_id` int(11) NOT NULL COMMENT 'ref:sys_permission.id',
    `created_at` datetime NOT NULL,
    UNIQUE KEY `uni_srp_sri_spi` (`sys_role_id`,`sys_permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_white_list` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `system_id` int(11) NOT NULL COMMENT 'ref:system.id',
    `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `email` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `type` enum('email','phone') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'email:檢查email欄位 phone:檢查phone欄位',
    `is_disable` tinyint(4) NOT NULL COMMENT '0:啟用 1:禁用',
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_swl_sid_phone_email_type` (`system_id`,`phone`,`email`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE `system`;
DROP TABLE `sys_account`;
DROP TABLE `sys_account_role`;
DROP TABLE `sys_permission`;
DROP TABLE `sys_purchase`;
DROP TABLE `sys_role`;
DROP TABLE `sys_role_permission`;
DROP TABLE `sys_white_list`;
