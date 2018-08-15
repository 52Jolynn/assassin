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
	unique key n(`name`)
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
	`status` varchar(4) not null comment '状态, N: 正常, D: 禁用',
	primary key(`id`),
	key n(`name`),
	key c(`club_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='俱乐部场地表';

##  俱乐部场地租用时段表
drop table if exists ground_rental;
create table if not exists ground_rental (
	`id` bigint not null auto_increment comment 'ID',
	`from_time` datetime(3) not null comment '可用起始时间',
	`to_time` datetime(3) not null comment '可用结束时间',
	`club_id` int not null comment '所属俱乐部id',
	`ground_id` int not null comment '场地id',
	`rent_amount` bigint not null default 0 comment '场租，单位分',
	`rel_rental_id` bigint null comment '关联的租用时段, 若有关联需打包租用，否则可单独租用',
	`create_time` datetime(3) not null comment '创建时间',
	`status` varchar(4) not null comment '状态, N: 正常, D: 禁用',
	primary key(`id`),
	key c(`club_id`),
	key g(`ground_id`),
	key cg(`club_id`, `ground_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='俱乐部场地租用时段表';

## 俱乐部场地租用记录表
drop table if exists ground_rental_record;
create table if not exists ground_rental_record (
	`id` bigint not null auto_increment comment 'ID',
	`rent_player_id` int not null comment '租用场地球员id',
	`from_rental_id` int not null comment '租用时段id',
	`to_rental_id` int not null comment '租用时段id',
	`from_time` datetime(3) not null comment '租用起始时间',
	`to_time` datetime(3) not null comment '租用结束时间',
	`rent_date` date not null comment '租用日期',
	`club_id` int not null comment '所属俱乐部id',
	`ground_id` int not null comment '场地id',
	`rent_amount` bigint not null default 0 comment '总场租，单位分',
	`create_time` datetime(3) not null comment '创建时间',
	`status` varchar(4) not null comment '状态, N: 正常, L: 已锁定, R: 已出租, T: 已转租',
	primary key(`id`),
	key c(`club_id`),
	key g(`ground_id`),
	key cg(`club_id`, `ground_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='俱乐部场地租用记录表';


## 球队表
drop table if exists team;
create table if not exists team (
	`id` int not null auto_increment comment 'ID',
	`name` varchar(32) not null comment '球队名称',
	`remark` varchar(128) null comment '描述',
	`captain_name` varchar(16) null comment '队长名称',
	`captain_mobile` varchar(16) null comment '队长联系电话',
	`manager_username` varchar(16) null comment '管理员用户名',
	`manager_passwd` varchar(128) null comment '管理员密码',
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
	key ct(`club_id`, `team_id`)
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
	key ct(`club_id`, `team_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='优惠券表';

## 球员表
drop table if exists player;
create table if not exists player (
	`id` int not null auto_increment comment 'ID',
	`username` varchar(16) null comment '用户名',
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
	`level` varchar(4) not null comment 'N: 普通队员(只能查看个人相关数据), S: 正式队员(可查看球队相关数据)', 
	primary key(`id`),
	key u(`username`),
	key name(`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='球员表';

## 球员数值评估表(自评+他评)
drop table if exists player_evaluation;
create table if not exists player_evaluation (
	`id` int not null auto_increment comment 'ID',
	`player_id` int not null comment '球员id',
	`team_id` int not null comment '球队id',
	`evaluate_player_id` int not null comment '作出评价的球员id',
	`fit` numeric not null comment '当前数值吻合度, 0%-100%',
	`create_time` datetime(3) not null comment '创建时间',
	primary key(`id`),
	key p(`player_id`),
	key t(`team_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='球员数值评估表';

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
	primary key(`id`),
	key p(`player_id`),
	key t(`team_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='球员与球队表';

## 球衣与球队表
drop table if exists jersey_of_team;
create table if not exists jersey_of_team (
	`id` int not null auto_increment comment 'ID',
	`team_id` int not null comment '球队id',
	`home_color` varchar(16) not null comment '主场球衣颜色',
	`away_color` varchar(16) null comment '客场球衣颜色',
	`create_time` datetime(3) not null comment '创建时间',
	`status` varchar(4) not null comment '状态, N: 正常, D: 禁用',
	primary key(`id`),
	key t(`team_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='球衣与球队表';

## 比赛表
drop table if exists game_of_match;
create table if not exists game_of_match (
	`id` bigint not null auto_increment comment 'ID',
	`home_team_id` int not null comment '主场球队id',
	`away_team_id` int not null comment '客场球队id',
	`club_id` int not null comment '俱乐部id，比赛场地',
	`ground_id` int not null comment '比赛场地类型id',
	`home_jersey_color` varchar(16) not null comment '主队球衣颜色',
	`away_jersey_color` varchar(16) not null comment '客队球衣颜色',
	`open_time` datetime(3) not null comment '开赛时间',
	`enroll_start_time` datetime(3) not null comment '开始报名时间',
	`enroll_end_time` datetime(3) not null comment '截止报名时间',
	`enroll_quota` int not null comment '报名人数上限',
	`rent_cost` bigint not null comment '总场租，单位元',
	`match_duration` int not null comment '比赛时长',
	`create_time` datetime(3) not null comment '创建时间',
	`status` varchar(4) not null comment '状态, O: 未开赛, C: 取消, P: 开赛进行中, E: 已结束',
	primary key(`id`),
	key ht(`home_team_id`),
	key c(`club_id`),
	key g(`ground_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='比赛表';

## 比赛报名表
drop table if exists enroll_of_match;
create table if not exists enroll_of_match (
	`id` bigint not null auto_increment comment 'ID',
	`match_id` bigint not null comment '比赛id',
	`player_id` int not null comment '球员id',
	`temporary_player` int not null default 0 comment '携带散兵数',
	`create_time` datetime(3) not null comment '创建时间',
	`status` varchar(4) not null comment '报名状态, F: 报名失败, S: 报名成功, C: 取消报名',
	primary key(`id`),
	key m(`match_id`),
	key p(`player_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='比赛报名表';

## 球队比赛统计表
drop table if exists team_stat_of_match;
create table if not exists team_stat_of_match (
	`id` bigint not null auto_increment comment 'ID',
	`match_id` bigint not null comment '比赛id',
	`type` varchar(8) not null comment '主队客队, home or away',
	`team_id` int not null comment '比赛id',
	`score` int not null default 0 comment '主队进球',
	`rent_amount` bigint not null comment '球队需承担的场租费用，单位分',
	`shot` int not null default 0 comment '射门数',
	`foul` int not null default 0 comment '犯规数',
	`free_kick` int not null default 0 comment '任意球数',
	`penalty_kick` int not null default 0 comment '任意球数',
	`offside` int not null default 0 comment '越位次数',
	`corner` int not null default 0 comment '角球数',
	`pass` int not null default 0 comment '传球数',
	`yellow_card` int not null default 0 comment '黄牌数',
	`red_card` int not null default 0 comment '红牌数',
	`create_time` datetime(3) not null comment '创建时间',
	primary key(`id`),
	key m(`match_id`),
	key t(`team_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='球队比赛统计表';

## 比赛球员统计表
drop table if exists player_stat_of_match;
create table if not exists player_stat_of_match (
	`id` bigint not null auto_increment comment 'ID',
	`match_id` bigint not null comment '比赛id',
	`player_id` int not null comment '球员id',
	`rent_amount` bigint not null comment '个人需承担的场租费用，单位分',
	`temporary_player_rent_amount` bigint not null default 0 comment '携带的散兵需承担的场租费用，单位分',
	`score` int not null default 0 comment '进球数',
	`shot` int not null default 0 comment '射门数',
	`assists` int not null default 0 comment '助攻数',
	`foul` int not null default 0 comment '犯规数',
	`break_through` int not null default 0 comment '过人数',
	`tackle` int not null default 0 comment '抢截数',
	`yellow_card` int not null default 0 comment '黄牌数',
	`red_card` int not null default 0 comment '红牌数',
	`create_time` datetime(3) not null comment '创建时间',
	`player_status` varchar(4) not null comment '到场状态, N: 到场, X: 未到场',
	`pay_by_sb` char(1) not null comment '是否由他人代交场租, Y or N',
	`pay_player_id` int null comment '代付场租球员id',
	primary key(`id`),
	key m(`match_id`),
	key p(`player_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='比赛球员统计表';

## 球队账目表
drop table if exists accounting_of_team;
create table if not exists accounting_of_team (
	`id` varchar(64) not null comment '流水ID',
	`ref_id` varchar(64) null comment '关联的流水ID',
	`match_id` bigint not null comment '比赛id',
	`team_id` int not null comment '球队id',
	`amount` bigint not null comment ' 金额，收入为正，支出为负',
	`remark` varchar(128) not null comment '备注',
	`before_balance` bigint not null comment '记账前余额',
	`after_balance` bigint not null comment '记账后余额',
	`type` varchar(4) not null comment '记账类型, SR(收入), ZZ(支出), CZ(冲正)',
	`bill_date` date not null comment '账目时间',
	`create_time` datetime(3) not null comment '创建时间',
	primary key(`id`),
	key m(`match_id`),
	key p(`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='球队账目表';

## 球员账目表
drop table if exists accounting_of_player;
create table if not exists accounting_of_player (
	`id` varchar(64) not null comment '流水ID',
	`ref_id` varchar(64) null comment '关联的流水ID',
	`match_id` bigint not null comment '比赛id',
	`team_id` int not null comment '球队id',
	`player_id` int not null comment '球员id',
	`amount` bigint not null comment ' 金额，充值为正，消费、提现为负',
	`remark` varchar(128) not null comment '备注',
	`before_balance` bigint not null comment '记账前余额',
	`after_balance` bigint not null comment '记账后余额',
	`bill_date` date not null comment '账目时间',
	`type` varchar(4) not null comment '记账类型, R(充值), C(消费), W(提现), CZ(冲正)',
	`create_time` datetime(3) not null comment '创建时间',
	primary key(`id`),
	key m(`match_id`),
	key p(`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='球员账目表';


```