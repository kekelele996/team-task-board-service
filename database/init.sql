-- Team Task Kanban MySQL bootstrap. The backend runs GORM AutoMigrate at startup;
-- this script only guarantees the database/user exist on first initialization.
CREATE DATABASE IF NOT EXISTS `gbkanban` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'gbkanban'@'%' IDENTIFIED BY 'gbkanban';
GRANT ALL PRIVILEGES ON `gbkanban`.* TO 'gbkanban'@'%';
FLUSH PRIVILEGES;
