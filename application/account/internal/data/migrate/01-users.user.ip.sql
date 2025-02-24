create table users1(
    id int8 primary key,
    name varchar(20)
);

DROP table users;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    mobile VARCHAR(11) UNIQUE NOT NULL,
    password VARCHAR(200) NOT NULL,
    nickname VARCHAR(25),
    birthday BIGINT,
    gender VARCHAR(16) DEFAULT 'male',
    role INT DEFAULT 1,
    created_at bigint,
    updated_at bigint,
    deleted_at bigint,
    is_deleted_at BOOLEAN
);
COMMENT ON COLUMN users.Mobile IS '手机号码, 用户唯一标识';
COMMENT ON COLUMN users.Birthday IS '出生日期';
COMMENT ON COLUMN users.Gender IS 'male: 男, female: 女';
COMMENT ON COLUMN users.Role IS '1:普通用户, 2:管理员';

