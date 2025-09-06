/*
 Navicat Premium Data Transfer

 Source Server         : 索引准生产
 Source Server Type    : MySQL
 Source Server Version : 80036 (8.0.36)
 Source Host           : 18.167.42.160:3306
 Source Schema         : robot_index

 Target Server Type    : MySQL
 Target Server Version : 80036 (8.0.36)
 File Encoding         : 65001

 Date: 06/09/2025 14:07:58
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for tele_index_record
-- ----------------------------
DROP TABLE IF EXISTS `tele_index_record`;
CREATE TABLE `tele_index_record` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '类型',
  `chat_id` bigint DEFAULT NULL COMMENT '聊天id',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '标题',
  `description` varchar(1500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '描述',
  `source` int DEFAULT '0' COMMENT '群组来源: 0-抓取,1-汇旺,2-采集',
  `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '标签',
  `keywords` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '关键字',
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '用户名',
  `link` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '链接',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '创建人',
  `bot_id` bigint DEFAULT NULL COMMENT '机器人',
  `members` bigint DEFAULT '0' COMMENT '总人数',
  `online` bigint DEFAULT NULL COMMENT '在线人数',
  `actives` bigint DEFAULT NULL COMMENT '活跃人数',
  `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '正常' COMMENT '状态',
  `isofficial` varchar(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT 'N' COMMENT '是否为官方认证,Y是N否',
  `isrepeat` int DEFAULT '0' COMMENT '是否重复: 0 否, 1 是',
  `offline_reason` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '下架原因',
  `offline_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '下架操作人',
  `likes` int DEFAULT '0' COMMENT '点赞数量',
  `weight` bigint DEFAULT NULL COMMENT '权重值',
  `dept_id` bigint DEFAULT NULL COMMENT '部门id',
  `language` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '语言',
  `avatar` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '头像',
  `picpath` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '图片路径',
  `update_status` int DEFAULT '0' COMMENT '更新状态',
  `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `start_date` datetime DEFAULT NULL,
  `end_date` datetime DEFAULT NULL,
  `category_id` bigint DEFAULT NULL COMMENT '分类id',
  `link_type` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '链接分类',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_status` (`status`) USING BTREE,
  KEY `idx_weight` (`weight`) USING BTREE,
  KEY `idx_link` (`link`) USING BTREE,
  KEY `idx_category` (`category_id`) USING BTREE,
  KEY `idx_chat_id` (`chat_id`) USING BTREE,
  KEY `idx_sort` (`weight` DESC,`members` DESC),
  FULLTEXT KEY `idx_title` (`title`) /*!50100 WITH PARSER `ngram` */ ,
  FULLTEXT KEY `idx_description` (`description`) /*!50100 WITH PARSER `ngram` */ ,
  FULLTEXT KEY `idx_title_description` (`title`,`description`) /*!50100 WITH PARSER `ngram` */ 
) ENGINE=InnoDB AUTO_INCREMENT=1921125144496320571 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='索引库信息表';

SET FOREIGN_KEY_CHECKS = 1;
