# 数据字典

```sql

CREATE DATABASE IF NOT EXISTS assassin DEFAULT CHARSET utf8 COLLATE utf8_general_ci;

use assassin;

##  线下俱乐部表
drop table if exists club;
create table if not exists club (
	`id` int not null auto_increment comment 'ID',
	`name` varchar(32) not null comment '名称(如爱奇、新安)',
	`remark` varchar(128) null comment '描述',
	`address` varchar(255) null comment '俱乐部详细地址',
	`tel` varchar(32) null comment '联系电话',
	`create_time` datetime(3) not null comment '创建时间',
	`status` varchar(4) not null comment '状态, N: 正常, D: 禁用',
	primary key(`id`),
	key n(`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='线下俱乐部表';

##  俱乐部场地表
drop table if exists ground;
create table if not exists ground (
	`id` int not null auto_increment comment 'ID',
	`name` varchar(32) not null comment '场地名称',
	`remark` varchar(128) null comment '描述',
	`type` varchar(16) not null comment '场地类型, 如七人场',
	`club_id` int not null comment '所属俱乐部id',
	`create_time` datetime(3) not null comment '创建时间',
	`status` varchar(4) not null comment '状态, N: 正常, U: 使用中, D: 禁用',
	primary key(`id`),
	key n(`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='俱乐部场地表';

## 球队表
drop table if exists team;
create table if not exists team (
	`id` int not null auto_increment comment 'ID',
	`name` varchar(32) not null comment '球队名称',
	`remark` varchar(128) null comment '描述',
	`captain_name` varchar(16) null comment '队长名称',
	`captain_mobile` varchar(16) null comment '队长联系电话',
	`manager_username` varchar(16) not null comment '管理员用户名',
	`manager_passwd` varchar(128) not null comment '管理员密码',
	`create_time` datetime(3) not null comment '创建时间',
	`status` varchar(4) not null comment '状态, N: 正常, D: 禁用',
	primary key(`id`),
	key n(`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='球队表';

## 球队与俱乐部表
drop table if exists team_of_club;
create table if not exists team_of_club (
	`id` int not null auto_increment comment 'ID',
	`club_id` int not null comment '俱乐部id',
	`team_id` int not null comment '球队id',
	`present_balance` bigint not null default 0 comment '球队余额(不包含优惠券)，单位分',
	`used_amount` bigint not null default 0 comment '已消费金额(不包含优惠券)，单位分',
	`join_time` datetime(3) null comment '加入时间',
	`create_time` datetime(3) not null comment '创建时间',
	primary key(`id`),
	key ct(`club_id``, `team_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='球队与俱乐部表';

## 优惠券表
drop table if exists coupon;
create table if not exists coupon (
	`id` int not null auto_increment comment 'ID',
	`club_id` int not null comment '俱乐部id',
	`team_id` int not null comment '球队id',
	`amount` bigint not null default 0 comment '优惠券金额，单位分',
	`least_amount` bigint not null default 0 comment '优惠券起用金额，单位分',
	`effective_time` datetime(3) not null comment '生效时间',
	`used_time` datetime(3) null comment '使用时间',
	`expire_time` datetime(3) null comment '过期时间',
	`create_time` datetime(3) not null comment '领取时间',
	`status` varchar(4) not null comment '状态, N: 正常, E: 已过期, U: 已使用',
	primary key(`id`),
	key ct(`club_id``, `team_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='优惠券表';

## 球员表
drop table if exists player;
create table if not exists player (
	`id` int not null auto_increment comment 'ID',
	`username` varchar(16) not null comment '用户名',
	`passwd` varchar(128) null comment '密码',
	`wx_open_id` varchar(128) null comment '微信openid',
	`name` varchar(32) not null comment '球员名称',
	`remark` varchar(128) null comment '描述',
	`mobile` varchar(16) null comment '联系电话',
	`pos` varchar(8) null comment '擅长位置, 多个以,号分隔',
	`height` numeric null comment '身高, 单位m',
	`age` int null comment '年龄',
	`pass_val` numeric null comment '传球',
	`shot_val` numeric null comment '射门',
	`strength_val` numeric null comment '力量',
	`dribble_val` numeric null comment '盘带',
	`speed_val` numeric null comment '速度',
	`tackle_val` numeric null comment '抢截',
	`head_val` numeric null comment '头球',
	`throwing_val` numeric null comment '手抛球',
	`reaction_val` numeric null comment '反应',
	`create_time` datetime(3) not null comment '创建时间',
	`status` varchar(4) not null comment '状态, N: 正常, E: 退出, D: 禁用',
	primary key(`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='球员表';

## 球员数值评估表(自评+他评)

## 球员与球队表
drop table if exists player_of_team;
create table if not exists player_of_team (
	`id` int not null auto_increment comment 'ID',
	`player_id` int not null comment '球员id',
	`team_id` int not null comment '球队id',
	`no` varchar(4) null comment '球员号码',
	`present_balance` bigint not null default 0 comment '余额，单位分',
	`used_amount` bigint not null default 0 comment '已消费金额，单位分',
	`join_time` datetime(3) null comment '加入时间',
	`create_time` datetime(3) not null comment '创建时间',
	primary key(`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='球员与球队表';

## 球衣与球队表
drop table if exists jersey_of_team;
create table if not exists jersey_of_team (
	`id` int not null auto_increment comment 'ID',
	`team_id` int not null comment '球队id',
	`home_color` varchar(16) not null comment '主场球衣颜色',
	`away_color` varchar(16) not null comment '客场球衣颜色',
	`create_time` datetime(3) not null comment '创建时间',
	primary key(`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='球衣与球队表';

## 比赛表
drop table if exists match;
create table if not exists match (
	`id` int not null auto_increment comment 'ID',
	`home_team_id` int not null comment '主场球队id',
	`club_id` int not null comment '俱乐部id，比赛场地',
	`ground_id` int not null comment '比赛场地类型id',
	`opponent` varchar(16) not null comment '对手',
	`home_jersey_id` varchar(16) not null comment '客队球衣颜色',
	`away_color` varchar(16) not null comment '客队球衣颜色',
	`open_time` datetime(3) not null comment '开赛时间',
	`enroll_start_time` datetime(3) not null comment '开始报名时间',
	`enroll_end_time` datetime(3) not null comment '截止报名时间',
	`enroll_quota` datetime(3) not null comment '报名人数上限',
	`rent_cost` bigint not null comment '总场租，单位元',
	`match_duration` int not null comment '比赛时长',
	`create_time` datetime(3) not null comment '创建时间',
	primary key(`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='比赛表';

## 比赛报名表

## 比赛统计表

## 球队账目表

## 球员账目表

```