ALTER TABLE `user_accounts`
DROP COLUMN `activation_token`,
DROP COLUMN `activation_otp`,
DROP COLUMN `activation_expires_at`;
