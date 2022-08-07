SET
    GLOBAL local_infile = 1;

CREATE TABLE IF NOT EXISTS `race` (
    `id` int(10) AUTO_INCREMENT NOT NULL PRIMARY KEY,
    `day_raceid` int(10) NOT NULL,
    `race_date` date NOT NULL,
    `start_time` timestamp NOT NULL,
    `end_time` timestamp,
    `temperature` float(5),
    `venue` varchar(50) NOT NULL,
    `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;