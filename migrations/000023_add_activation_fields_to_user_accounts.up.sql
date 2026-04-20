ALTER TABLE `user_accounts` 
ADD COLUMN `activation_token` VARCHAR(255) NULL,
ADD COLUMN `activation_otp` VARCHAR(6) NULL,
ADD COLUMN `activation_expires_at` DATETIME NULL;
