-- Dependencies (Referenced Tables)
CREATE TABLE IF NOT EXISTS access_control_servers (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name VARCHAR(255) NOT NULL,
host_address VARCHAR(255),
username VARCHAR(255),
password VARCHAR(255),
access_token VARCHAR(255),
api_token VARCHAR(255),
status VARCHAR(255),
last_sync_at TIMESTAMP WITH TIME ZONE,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS access_control_groups (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name VARCHAR(255) NOT NULL,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS access_control_rules (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name VARCHAR(255) NOT NULL,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS attendances (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name VARCHAR(255) NOT NULL,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS user_permissions (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
people_permission BOOLEAN DEFAULT FALSE,
device_permission BOOLEAN DEFAULT FALSE,
rule_permission BOOLEAN DEFAULT FALSE,
time_attendance_permission BOOLEAN DEFAULT FALSE,
report_permission BOOLEAN DEFAULT FALSE,
notification_permission BOOLEAN DEFAULT FALSE,
system_log_permission BOOLEAN DEFAULT FALSE,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS register_forms (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name VARCHAR(255) NOT NULL,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS people (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
title VARCHAR(255),
first_name VARCHAR(255),
middle_name VARCHAR(255),
last_name VARCHAR(255),
gender VARCHAR(255),
date_of_birth DATE,
person_type VARCHAR(50),
person_id VARCHAR(255),
company VARCHAR(255),
department VARCHAR(255),
job_position VARCHAR(255),
address TEXT,
mobile_number VARCHAR(20),
email VARCHAR(255),
face_image_path VARCHAR(255),
is_verified BOOLEAN DEFAULT FALSE,
active_at TIMESTAMP WITH TIME ZONE,
expire_at TIMESTAMP WITH TIME ZONE,
rule_id UUID,
time_attendance_id UUID,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (first_name, last_name, person_id)
);

-- Referencing Tables (Level 1)
CREATE TABLE IF NOT EXISTS access_control_devices (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name VARCHAR(255) NOT NULL,
type VARCHAR(255),
host_address VARCHAR(255),
username VARCHAR(255),
password VARCHAR(255),
access_token VARCHAR(255),
api_token VARCHAR(255),
server_id UUID REFERENCES access_control_servers(id),
record_scan BOOLEAN DEFAULT FALSE,
record_attendance BOOLEAN DEFAULT FALSE,
allow_clock_in BOOLEAN DEFAULT FALSE,
allow_clock_out BOOLEAN DEFAULT FALSE,
status VARCHAR(255),
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS access_control_group_schedules (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
access_control_group_id UUID NOT NULL REFERENCES access_control_groups(id) ON DELETE CASCADE,
day_of_week INTEGER,
date DATE,
start_time TIME,
end_time TIME,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS attendance_schedules (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
attendance_id UUID NOT NULL REFERENCES attendances(id) ON DELETE CASCADE,
day_of_week INTEGER,
date DATE,
start_time TIME,
end_time TIME,
early_in_minutes INTEGER,
late_in_minutes INTEGER,
early_out_minutes INTEGER,
late_out_minutes INTEGER,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS person_cards (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
card_number VARCHAR(255) UNIQUE NOT NULL,
person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
active_at TIMESTAMP WITH TIME ZONE,
expire_at TIMESTAMP WITH TIME ZONE,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (card_number)
);

CREATE TABLE IF NOT EXISTS person_license_plates (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
license_plate_number VARCHAR(255) UNIQUE NOT NULL,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (license_plate_number, person_id)
);

CREATE TABLE IF NOT EXISTS register_form_fields (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name VARCHAR(255) NOT NULL,
register_form_id UUID NOT NULL REFERENCES register_forms(id) ON DELETE CASCADE,
field_order INTEGER,
field_type VARCHAR(50),
input_type VARCHAR(50),
placeholder VARCHAR(255),
label VARCHAR(255),
help_text TEXT,
is_required BOOLEAN DEFAULT FALSE,
default_value TEXT,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (name, register_form_id)
);

CREATE TABLE IF NOT EXISTS users (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
username VARCHAR(255) UNIQUE NOT NULL,
password_hash VARCHAR(255) NOT NULL,
permission_id UUID NOT NULL REFERENCES user_permissions(id) ON DELETE CASCADE,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (username)
);

-- Referencing Tables (Level 2)
CREATE TABLE IF NOT EXISTS access_control_group_devices (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
access_control_group_id UUID NOT NULL REFERENCES access_control_groups(id) ON DELETE CASCADE,
access_control_device_id UUID NOT NULL REFERENCES access_control_devices(id) ON DELETE CASCADE,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (access_control_group_id, access_control_device_id)
);

CREATE TABLE IF NOT EXISTS access_control_rule_groups (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
access_control_group_id UUID NOT NULL REFERENCES access_control_groups(id) ON DELETE CASCADE,
access_control_device_id UUID NOT NULL REFERENCES access_control_devices(id) ON DELETE CASCADE,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (access_control_group_id, access_control_device_id)
);

CREATE TABLE IF NOT EXISTS access_records (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
access_control_device_id UUID REFERENCES access_control_devices(id) ON DELETE CASCADE,
type VARCHAR(255),
result VARCHAR(255),
access_time TIMESTAMP WITH TIME ZONE,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS register_form_field_answers (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
register_form_field_id UUID NOT NULL REFERENCES register_form_fields(id) ON DELETE CASCADE,
register_form_id UUID NOT NULL REFERENCES register_forms(id) ON DELETE CASCADE,
person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
answer_value TEXT,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE,
UNIQUE (register_form_field_id, register_form_id, person_id)
);

-- Referencing Tables (Level 3)
CREATE TABLE IF NOT EXISTS attendance_records (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
attendance_schedule_id UUID NOT NULL REFERENCES attendance_schedules(id) ON DELETE CASCADE,
access_record_id UUID NOT NULL REFERENCES access_records(id) ON DELETE CASCADE,
date DATE,
status VARCHAR(255),
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
deleted_at TIMESTAMP WITH TIME ZONE
);