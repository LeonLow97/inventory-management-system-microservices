#!/bin/bash

sudo dnf install -y postgresql15

USER="ec2-user"

# Fetching Bastion host IP and key name dynamically
BASTION_PUBLIC_IP_AZ1=$(aws ec2 describe-instances --filters "Name=tag:Name,Values=BastionHost" --query "Reservations[0].Instances[0].PublicIpAddress" --output text)
EC2_BASTION_HOST_KEY_NAME_AZ1=IMS_BASTION_AZ1

# Fetching RDS configuration from AWS SSM
RDS_ENDPOINT=$(aws ssm get-parameters --names "/ims/db/hostname" --query "Parameters[0].Value" --output text)
DB_MASTER_USERNAME=$(aws ssm get-parameters --names "/ims/db/master-username" --query "Parameters[0].Value" --output text)
DB_MASTER_PASSWORD=$(aws ssm get-parameters --names "/ims/db/master-password" --with-decryption --query "Parameters[0].Value" --output text)
IMS_DB_NAME=$(aws ssm get-parameters --names "/ims/db/db-name" --query "Parameters[0].Value" --output text)

# Install PostgreSQL client
sudo dnf install -y postgresql15

# Create imsdb database
echo "Creating $IMS_DB_NAME database..."
PGPASSWORD=$DB_MASTER_PASSWORD psql -h $RDS_ENDPOINT -U $DB_MASTER_USERNAME -d postgres -c "CREATE DATABASE $IMS_DB_NAME;" || { echo "❌ Failed to create database $IMS_DB_NAME"; exit 1; }
echo "✅ Successfully created database $IMS_DB_NAME."

# Grant privileges to the master user
echo "Granting privileges to $DB_MASTER_USERNAME on $IMS_DB_NAME..."
PGPASSWORD=$DB_MASTER_PASSWORD psql -h $RDS_ENDPOINT -U $DB_MASTER_USERNAME -d $IMS_DB_NAME -c "
  -- Grant privileges on $IMS_DB_NAME to the master user
  GRANT ALL PRIVILEGES ON DATABASE $IMS_DB_NAME TO $DB_MASTER_USERNAME;

  -- Grant privileges on public schema to $DB_MASTER_USERNAME
  GRANT ALL PRIVILEGES ON SCHEMA public TO $DB_MASTER_USERNAME;

  -- Grant SELECT, INSERT, UPDATE, DELETE on all tables in the public schema
  GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO $DB_MASTER_USERNAME;

  -- Set default privileges for new tables and sequences in the public schema
  ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $DB_MASTER_USERNAME;
  ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO $DB_MASTER_USERNAME;
" || { echo "❌ Failed to grant privileges to $DB_MASTER_USERNAME"; exit 1; }
echo "✅ Successfully granted privileges to $DB_MASTER_USERNAME."

# Copy the SQL initialization file to the Bastion host
echo "Copying SQL initialization file to Bastion host..."
scp -i "$EC2_BASTION_HOST_KEY_NAME_AZ1" "./init-db/init-authentication-db.sql" "$USER@$BASTION_PUBLIC_IP_AZ1:~" || { echo "❌ Failed to copy SQL file to Bastion host"; exit 1; }
echo "✅ Successfully copied SQL file."

# Run the SQL script to initialize the database
echo "Running the SQL script to initialize database tables..."
PGPASSWORD=$DB_MASTER_PASSWORD psql -h $RDS_ENDPOINT -U $DB_MASTER_USERNAME -d $IMS_DB_NAME -f ./init-db/init-authentication-db.sql || { echo "❌ Failed to execute SQL script"; exit 1; }
echo "✅ Successfully created tables in Authentication Postgres Database."