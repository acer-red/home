-- Don't modify this.

CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),
    
    username        VARCHAR(20) UNIQUE,
    email           VARCHAR(255) UNIQUE,
    password        TEXT,
    public_key      BYTEA,
    nickname        TEXT,
    avatar_name     TEXT,
    avatar_url      TEXT,


    created_at      TIMESTAMP WITH TIME ZONE,
    updated_at      TIMESTAMP WITH TIME ZONE,
    deleted_at      TIMESTAMP WITH TIME ZONE
);
COMMENT ON TABLE users IS '统一用户表';
COMMENT ON COLUMN users.id IS '主键,uuidv7格式,全局用户ID';
COMMENT ON COLUMN users.public_key IS '用户公钥，用于加密敏感数据';



CREATE TABLE IF NOT EXISTS product (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),

    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    category        VARCHAR(255),

    created_at      TIMESTAMP WITH TIME ZONE
);
COMMENT ON TABLE product IS '产品注册表';
COMMENT ON COLUMN product.id IS '主键,uuidv7格式,同一用户ID下的不同产品ID,产品ID';
COMMENT ON COLUMN product.user_id IS '关联的用户ID';
COMMENT ON COLUMN product.category IS '产品类别';



CREATE TABLE IF NOT EXISTS device (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),
    product_id      UUID REFERENCES product(id) ON DELETE CASCADE,

    created_at      TIMESTAMP WITH TIME ZONE,
    updated_at      TIMESTAMP WITH TIME ZONE,
    deleted_at      TIMESTAMP WITH TIME ZONE
);
COMMENT ON TABLE device IS '设备表';
COMMENT ON COLUMN device.id IS '主键,uuidv7格式,同一产品ID下的不同设备ID,设备ID';
COMMENT ON COLUMN device.product_id IS '关联的产品ID';



CREATE TABLE IF NOT EXISTS api (
    id              BIGSERIAL PRIMARY KEY,
    product_id      UUID REFERENCES product(id) ON DELETE CASCADE,
    
    api_key         VARCHAR(255) UNIQUE,
    expires_at      TIMESTAMP WITH TIME ZONE,
    last_used_at    TIMESTAMP WITH TIME ZONE,
    used_times      INTEGER,
    
    created_at      TIMESTAMP WITH TIME ZONE,
    updated_at      TIMESTAMP WITH TIME ZONE,
    deleted_at      TIMESTAMP WITH TIME ZONE
);
COMMENT ON TABLE api IS 'API 密钥表';
COMMENT ON COLUMN api.api_key IS '唯一的 API 密钥';
COMMENT ON COLUMN api.product_id IS '关联的产品 ID';



CREATE TABLE IF NOT EXISTS files (
    id              BIGSERIAL PRIMARY KEY,
    user_id         UUID REFERENCES users(id)  ON DELETE CASCADE,

    name            VARCHAR(255) UNIQUE,
    data            BYTEA,
    category        TEXT,
    mime_type       TEXT,
    size            BIGINT,
    metadata        TEXT,  

    created_at      TIMESTAMP WITH TIME ZONE,
    updated_at      TIMESTAMP WITH TIME ZONE,
    deleted_at      TIMESTAMP WITH TIME ZONE
);
COMMENT ON TABLE files IS '文件存储表';
COMMENT ON COLUMN files.id IS '主键,';
COMMENT ON COLUMN files.user_id IS '文件所属用户 ID';
COMMENT ON COLUMN files.name IS '文件唯一名称';
COMMENT ON COLUMN files.data IS '文件的二进制数据';



CREATE TABLE IF NOT EXISTS feedback (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),
    product_id      UUID REFERENCES product(id),

    fb_type         INTEGER,
    title           TEXT,
    content         TEXT,
    is_public       BOOLEAN,
    device_file     TEXT,
    images          JSONB,

    created_at      TIMESTAMP WITH TIME ZONE,
    updated_at      TIMESTAMP WITH TIME ZONE,
    deleted_at      TIMESTAMP WITH TIME ZONE
);
COMMENT ON TABLE feedback IS '反馈表';
COMMENT ON COLUMN feedback.id IS '主键,uuidv7格式';
COMMENT ON COLUMN feedback.product_id IS '提交反馈的产品 ID';
COMMENT ON COLUMN feedback.images IS '关联的图片文件名列表，存储为 JSON 数组';
COMMENT ON COLUMN feedback.device_file IS '关联的设备文件名';


