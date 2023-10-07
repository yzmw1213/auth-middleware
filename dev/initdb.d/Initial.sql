DROP SCHEMA IF EXISTS user_db;
DROP SCHEMA IF EXISTS client_db;
CREATE SCHEMA user_db;
CREATE SCHEMA client_db;

-- 既存の場合、CREATEできない
CREATE USER IF NOT EXISTS 'yzmw1213'@'%' IDENTIFIED BY 'fga%45ng2eBj9d';
GRANT ALL ON user_db.* TO 'yzmw1213'@'%';
GRANT ALL ON client_db.* TO 'yzmw1213'@'%';

DROP TABLE IF EXISTS user_db.users;
DROP TABLE IF EXISTS client_db.client_users;

CREATE TABLE user_db.users
(
    user_id int(255) NOT NULL AUTO_INCREMENT COMMENT 'ユーザーID',
    name varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '名前',
    email varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'メールアドレス',
    firebase_uid varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'uid',
    delete_flag tinyint(1) NOT NULL DEFAULT '0' COMMENT '削除フラグ',
    create_user_id int(255) NOT NULL DEFAULT '0' COMMENT '作成ユーザーID',
    update_user_id int(255) NOT NULL DEFAULT '0' COMMENT '更新ユーザーID',
    created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `uk_uu_1` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ユーザー';

CREATE TABLE client_db.client_users
(
    client_user_id int(255) NOT NULL AUTO_INCREMENT COMMENT '利用ユーザーID',
    user_id int(255) NOT NULL DEFAULT '0' COMMENT 'ユーザーID',
    delete_flag tinyint(1) NOT NULL DEFAULT '0' COMMENT '削除フラグ',
    create_user_id int(255) NOT NULL DEFAULT '0' COMMENT '作成ユーザーID',
    update_user_id int(255) NOT NULL DEFAULT '0' COMMENT '更新ユーザーID',
    created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    PRIMARY KEY (`client_user_id`),
    FOREIGN KEY `fk_cu_1` (`user_id`) REFERENCES `user_db` . `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='利用ユーザー';

INSERT INTO user_db.users (name, email, firebase_uid) VALUES ('admin user', 'admin@example.co.jp', 'WzEyGeAl5BRcn3pcHACdTSRopfC3');
INSERT INTO user_db.users (name, email, firebase_uid) VALUES ('default user', 'default@example.co.jp', 'Y243jnTMQSOz7Pjkfgha5Vffcpl2');
INSERT INTO client_db.client_users (user_id) VALUES ((SELECT user_id FROM user_db.users WHERE firebase_uid='Y243jnTMQSOz7Pjkfgha5Vffcpl2'));
