create database video_server;
use video_server;

create table users(
id int UNSIGNED auto_increment  primary key,
user_name varchar(64) unique key not null,
pwd varchar(150) not null,
email varchar(64) unique key not null
);

create table video_info(
id int unsigned primary key auto_increment,
vid text not null,
author_name varchar(64)  not null,
title varchar(30) not null,
sub_tag tinyint unsigned not null,
parent_tag tinyint unsigned not null,
create_time datetime default current_timestamp
);



create table video_type(
type_id smallint unsigned not null primary key,
type_name varchar(20) not null
);


create table sessions(
session_id varchar(64) not null primary key,
expire varchar(20) not null,
user_name varchar(64) not null unique key
);


create table comments(
id int unsigned primary key auto_increment ,
video_id varchar(30) not null ,
user_name varchar(64) not null ,
content text not null ,
ctime datetime default current_timestamp
);


INSERT INTO video_type (type_id,type_name) VALUES 
(1,'电影')
,(2,'电视剧')
,(3,'动漫')
,(4,'综艺')
,(5,'短视频')
,(21,'沙雕>搞笑>动作>科幻')
,(22,'科幻')
,(23,'喜剧')
,(24,'动作')
,(25,'爱情')
;
INSERT INTO video_type (type_id,type_name) VALUES 
(26,'恐怖')
,(27,'国产')
,(28,'港台')
,(29,'日韩')
,(30,'欧美')
,(31,'沙雕')
;