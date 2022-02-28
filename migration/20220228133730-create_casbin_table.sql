
-- +migrate Up
CREATE TABLE `casbin_rule` (
    `p_type` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `v0` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `v1` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `v2` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `v3` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `v4` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `v5` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE `casbin_rule`;
