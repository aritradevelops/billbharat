-- Create System Admin
INSERT INTO users (
  id,
  human_id,
  name,
  email,
  phone,
  email_verified,
  phone_verified,
  created_at,
  created_by
) VALUES (
  '00000000-0000-0000-0000-000000000000',
  'system_admin',
  'System Admin',
  'systemadmin@billbharat.com',
  '919999999999',
  true,
  true,
  now(),
  '00000000-0000-0000-0000-000000000000'
);
