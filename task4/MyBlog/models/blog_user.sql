-- CREATE TABLE `myblogdb`.`user` (
--                                    `id` INT NOT NULL AUTO_INCREMENT,
--                                    `user_id` INT NOT NULL,
--                                    `username` VARCHAR(45) NOT NULL,
--                                    `password` VARCHAR(45) NOT NULL,
--                                    PRIMARY KEY (`id`))
--     ENGINE = InnoDB
-- DEFAULT CHARACTER SET = utf8mb4;



create table user
(
    id          bigint not null auto_increment primary key,
    user_id     bigint                              not null,
    username    varchar(64)                         not null,
    password    varchar(64)                         not null,
    email       varchar(64)                         null,
    gender      tinyint   default 0                 not null,
    create_time timestamp default CURRENT_TIMESTAMP null,
    update_time timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint idx_user_id
        unique (user_id),
    constraint idx_username
        unique (username)
)
    collate = utf8mb4_general_ci;

INSERT INTO myblogdb.user (id, user_id, username, password, email, gender, create_time, update_time) VALUES (1, 28018727488323585, 'tony', '313233343536639a9119599647d841b1bef6ce5ea293', null, 0, '2025-11-26 07:01:03', '2025-11-26 07:01:03');
INSERT INTO myblogdb.user (id, user_id, username, password, email, gender, create_time, update_time) VALUES (2, 4183532125556736, 'jack', '313233639a9119599647d841b1bef6ce5ea293', null, 0, '2025-11-26 13:03:51', '2025-11-26 13:03:51');