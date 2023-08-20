DROP SCHEMA IF EXISTS user_db;
CREATE SCHEMA user_db;

CREATE USER 'yzmw1213'@'%' IDENTIFIED BY 'fga%45ng2eBj9d';
GRANT ALL ON user_db.* TO 'yzmw1213'@'%';

DROP TABLE IF EXISTS user_db.users;

CREATE TABLE user_db.users
(
    user_id int(255) NOT NULL AUTO_INCREMENT COMMENT 'ユーザーID',
    name varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '名前',
    email varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'メールアドレス',
    firebase_uid varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'uid',
    delete_flag tinyint(1) NOT NULL DEFAULT '0' COMMENT '削除フラグ',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `uk_uu` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ユーザー';

INSERT INTO user_db.users (name, email, firebase_uid) VALUES ('default@example.co.jp', 'default user', 'Y243jnTMQSOz7Pjkfgha5Vffcpl2');