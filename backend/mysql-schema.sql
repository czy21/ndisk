create database if not exists cloud_disk_sync default charset utf8mb4 collate utf8mb4_0900_ai_ci;
use cloud_disk_sync;
CREATE TABLE provider_account (
  id       int NOT NULL AUTO_INCREMENT,
  username varchar(50) NULL,
  password varchar(50) NULL,
  token    varchar(50) NULL,
  kind     varchar(50) NOT NULL,
  create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  create_user int NULL,
  update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  update_user int NULL,
  deleted bit(1) NOT NULL DEFAULT b'0',
  PRIMARY KEY (`id`)
);
CREATE TABLE provider_folder (
 id int NOT NULL AUTO_INCREMENT,
 name        varchar(50) NOT NULL,
 remote_name varchar(50) NOT NULL,
 provider_account_id int NOT NULL,
 create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
 create_user int NULL,
 update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 update_user int NULL,
 deleted bit(1) NOT NULL DEFAULT b'0',
 PRIMARY KEY (`id`),
 INDEX `idx_ProviderAccountId`(`provider_account_id`),
 CONSTRAINT `fk_ProviderFolder_ProviderAccount_ProviderAccountId`
 FOREIGN KEY (`provider_account_id`)
 REFERENCES `provider_account` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
);