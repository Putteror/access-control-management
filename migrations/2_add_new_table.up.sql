CREATE TABLE IF NOT EXISTS device (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host_address VARCHAR(255),
    name VARCHAR(255),
    username VARCHAR(255),
    password VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
)