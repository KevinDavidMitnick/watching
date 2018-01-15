/*
Navicat MySQL Data Transfer

Source Server         : 10.21.1.214
Source Server Version : 50552
Source Host           : 10.21.1.214:3306
Source Database       : ops-cmdb

Target Server Type    : MYSQL
Target Server Version : 50552
File Encoding         : 65001

Date: 2017-12-25 10:32:15
*/


CREATE DATABASE IF NOT EXISTS `ops-cmdb` DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci;
USE ops-cmdb;

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for account_emailaddress
-- ----------------------------
CREATE TABLE IF NOT EXISTS `account_emailaddress` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(254) NOT NULL,
  `verified` tinyint(1) NOT NULL,
  `primary` tinyint(1) NOT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `account_emailaddress_user_id_2c513194_fk_auth_user_id` (`user_id`),
  CONSTRAINT `account_emailaddress_user_id_2c513194_fk_auth_user_id` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of account_emailaddress
-- ----------------------------

-- ----------------------------
-- Table structure for account_emailconfirmation
-- ----------------------------
CREATE TABLE IF NOT EXISTS `account_emailconfirmation` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created` datetime NOT NULL,
  `sent` datetime DEFAULT NULL,
  `key` varchar(64) NOT NULL,
  `email_address_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`),
  KEY `account_emailconfirm_email_address_id_5b7f8c58_fk_account_e` (`email_address_id`),
  CONSTRAINT `account_emailconfirm_email_address_id_5b7f8c58_fk_account_e` FOREIGN KEY (`email_address_id`) REFERENCES `account_emailaddress` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of account_emailconfirmation
-- ----------------------------

-- ----------------------------
-- Table structure for auth_group
-- ----------------------------
CREATE TABLE IF NOT EXISTS `auth_group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(80) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of auth_group
-- ----------------------------

-- ----------------------------
-- Table structure for auth_group_permissions
-- ----------------------------
CREATE TABLE IF NOT EXISTS `auth_group_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_id` int(11) NOT NULL,
  `permission_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth_group_permissions_group_id_permission_id_0cd325b0_uniq` (`group_id`,`permission_id`),
  KEY `auth_group_permissio_permission_id_84c5c92e_fk_auth_perm` (`permission_id`),
  CONSTRAINT `auth_group_permissions_group_id_b120cbf9_fk_auth_group_id` FOREIGN KEY (`group_id`) REFERENCES `auth_group` (`id`),
  CONSTRAINT `auth_group_permissio_permission_id_84c5c92e_fk_auth_perm` FOREIGN KEY (`permission_id`) REFERENCES `auth_permission` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of auth_group_permissions
-- ----------------------------

-- ----------------------------
-- Table structure for auth_permission
-- ----------------------------
CREATE TABLE IF NOT EXISTS `auth_permission` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `content_type_id` int(11) NOT NULL,
  `codename` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth_permission_content_type_id_codename_01ab375a_uniq` (`content_type_id`,`codename`),
  CONSTRAINT `auth_permission_content_type_id_2f476e4b_fk_django_co` FOREIGN KEY (`content_type_id`) REFERENCES `django_content_type` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=73 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of auth_permission
-- ----------------------------
INSERT INTO `auth_permission` VALUES ('1', 'Can add log entry', '1', 'add_logentry');
INSERT INTO `auth_permission` VALUES ('2', 'Can change log entry', '1', 'change_logentry');
INSERT INTO `auth_permission` VALUES ('3', 'Can delete log entry', '1', 'delete_logentry');
INSERT INTO `auth_permission` VALUES ('4', 'Can add permission', '2', 'add_permission');
INSERT INTO `auth_permission` VALUES ('5', 'Can change permission', '2', 'change_permission');
INSERT INTO `auth_permission` VALUES ('6', 'Can delete permission', '2', 'delete_permission');
INSERT INTO `auth_permission` VALUES ('7', 'Can add group', '3', 'add_group');
INSERT INTO `auth_permission` VALUES ('8', 'Can change group', '3', 'change_group');
INSERT INTO `auth_permission` VALUES ('9', 'Can delete group', '3', 'delete_group');
INSERT INTO `auth_permission` VALUES ('10', 'Can add user', '4', 'add_user');
INSERT INTO `auth_permission` VALUES ('11', 'Can change user', '4', 'change_user');
INSERT INTO `auth_permission` VALUES ('12', 'Can delete user', '4', 'delete_user');
INSERT INTO `auth_permission` VALUES ('13', 'Can add content type', '5', 'add_contenttype');
INSERT INTO `auth_permission` VALUES ('14', 'Can change content type', '5', 'change_contenttype');
INSERT INTO `auth_permission` VALUES ('15', 'Can delete content type', '5', 'delete_contenttype');
INSERT INTO `auth_permission` VALUES ('16', 'Can add session', '6', 'add_session');
INSERT INTO `auth_permission` VALUES ('17', 'Can change session', '6', 'change_session');
INSERT INTO `auth_permission` VALUES ('18', 'Can delete session', '6', 'delete_session');
INSERT INTO `auth_permission` VALUES ('19', 'Can add Token', '7', 'add_token');
INSERT INTO `auth_permission` VALUES ('20', 'Can change Token', '7', 'change_token');
INSERT INTO `auth_permission` VALUES ('21', 'Can delete Token', '7', 'delete_token');
INSERT INTO `auth_permission` VALUES ('22', 'Can add ci_attr_option', '8', 'add_ci_attr_option');
INSERT INTO `auth_permission` VALUES ('23', 'Can change ci_attr_option', '8', 'change_ci_attr_option');
INSERT INTO `auth_permission` VALUES ('24', 'Can delete ci_attr_option', '8', 'delete_ci_attr_option');
INSERT INTO `auth_permission` VALUES ('25', 'Can add ci_log_follower', '9', 'add_ci_log_follower');
INSERT INTO `auth_permission` VALUES ('26', 'Can change ci_log_follower', '9', 'change_ci_log_follower');
INSERT INTO `auth_permission` VALUES ('27', 'Can delete ci_log_follower', '9', 'delete_ci_log_follower');
INSERT INTO `auth_permission` VALUES ('28', 'Can add ci_zone', '10', 'add_ci_zone');
INSERT INTO `auth_permission` VALUES ('29', 'Can change ci_zone', '10', 'change_ci_zone');
INSERT INTO `auth_permission` VALUES ('30', 'Can delete ci_zone', '10', 'delete_ci_zone');
INSERT INTO `auth_permission` VALUES ('31', 'Can add ci_log', '11', 'add_ci_log');
INSERT INTO `auth_permission` VALUES ('32', 'Can change ci_log', '11', 'change_ci_log');
INSERT INTO `auth_permission` VALUES ('33', 'Can delete ci_log', '11', 'delete_ci_log');
INSERT INTO `auth_permission` VALUES ('34', 'Can add ci_interface', '12', 'add_ci_interface');
INSERT INTO `auth_permission` VALUES ('35', 'Can change ci_interface', '12', 'change_ci_interface');
INSERT INTO `auth_permission` VALUES ('36', 'Can delete ci_interface', '12', 'delete_ci_interface');
INSERT INTO `auth_permission` VALUES ('37', 'Can add ci_type', '13', 'add_ci_type');
INSERT INTO `auth_permission` VALUES ('38', 'Can change ci_type', '13', 'change_ci_type');
INSERT INTO `auth_permission` VALUES ('39', 'Can delete ci_type', '13', 'delete_ci_type');
INSERT INTO `auth_permission` VALUES ('40', 'Can add ci_attr', '14', 'add_ci_attr');
INSERT INTO `auth_permission` VALUES ('41', 'Can change ci_attr', '14', 'change_ci_attr');
INSERT INTO `auth_permission` VALUES ('42', 'Can delete ci_attr', '14', 'delete_ci_attr');
INSERT INTO `auth_permission` VALUES ('43', 'Can add ci_object_phy', '15', 'add_ci_object_phy');
INSERT INTO `auth_permission` VALUES ('44', 'Can change ci_object_phy', '15', 'change_ci_object_phy');
INSERT INTO `auth_permission` VALUES ('45', 'Can delete ci_object_phy', '15', 'delete_ci_object_phy');
INSERT INTO `auth_permission` VALUES ('46', 'Can add ci_object_virtual', '16', 'add_ci_object_virtual');
INSERT INTO `auth_permission` VALUES ('47', 'Can change ci_object_virtual', '16', 'change_ci_object_virtual');
INSERT INTO `auth_permission` VALUES ('48', 'Can delete ci_object_virtual', '16', 'delete_ci_object_virtual');
INSERT INTO `auth_permission` VALUES ('49', 'Can add ci_attr_keylist', '17', 'add_ci_attr_keylist');
INSERT INTO `auth_permission` VALUES ('50', 'Can change ci_attr_keylist', '17', 'change_ci_attr_keylist');
INSERT INTO `auth_permission` VALUES ('51', 'Can delete ci_attr_keylist', '17', 'delete_ci_attr_keylist');
INSERT INTO `auth_permission` VALUES ('52', 'Can add ci_object', '18', 'add_ci_object');
INSERT INTO `auth_permission` VALUES ('53', 'Can change ci_object', '18', 'change_ci_object');
INSERT INTO `auth_permission` VALUES ('54', 'Can delete ci_object', '18', 'delete_ci_object');
INSERT INTO `auth_permission` VALUES ('55', 'Can add ci_datacenter', '19', 'add_ci_datacenter');
INSERT INTO `auth_permission` VALUES ('56', 'Can change ci_datacenter', '19', 'change_ci_datacenter');
INSERT INTO `auth_permission` VALUES ('57', 'Can delete ci_datacenter', '19', 'delete_ci_datacenter');
INSERT INTO `auth_permission` VALUES ('58', 'Can add type_link_zone', '20', 'add_type_link_zone');
INSERT INTO `auth_permission` VALUES ('59', 'Can change type_link_zone', '20', 'change_type_link_zone');
INSERT INTO `auth_permission` VALUES ('60', 'Can delete type_link_zone', '20', 'delete_type_link_zone');
INSERT INTO `auth_permission` VALUES ('61', 'Can add type_link_attr_key', '21', 'add_type_link_attr_key');
INSERT INTO `auth_permission` VALUES ('62', 'Can change type_link_attr_key', '21', 'change_type_link_attr_key');
INSERT INTO `auth_permission` VALUES ('63', 'Can delete type_link_attr_key', '21', 'delete_type_link_attr_key');
INSERT INTO `auth_permission` VALUES ('64', 'Can add tableview', '22', 'add_tableview');
INSERT INTO `auth_permission` VALUES ('65', 'Can change tableview', '22', 'change_tableview');
INSERT INTO `auth_permission` VALUES ('66', 'Can delete tableview', '22', 'delete_tableview');
INSERT INTO `auth_permission` VALUES ('67', 'Can add table_item', '23', 'add_table_item');
INSERT INTO `auth_permission` VALUES ('68', 'Can change table_item', '23', 'change_table_item');
INSERT INTO `auth_permission` VALUES ('69', 'Can delete table_item', '23', 'delete_table_item');
INSERT INTO `auth_permission` VALUES ('70', 'Can add table_view', '24', 'add_table_view');
INSERT INTO `auth_permission` VALUES ('71', 'Can change table_view', '24', 'change_table_view');
INSERT INTO `auth_permission` VALUES ('72', 'Can delete table_view', '24', 'delete_table_view');

-- ----------------------------
-- Table structure for auth_user
-- ----------------------------
CREATE TABLE IF NOT EXISTS `auth_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `password` varchar(128) NOT NULL,
  `last_login` datetime DEFAULT NULL,
  `is_superuser` tinyint(1) NOT NULL,
  `username` varchar(150) NOT NULL,
  `first_name` varchar(30) NOT NULL,
  `last_name` varchar(30) NOT NULL,
  `email` varchar(254) NOT NULL,
  `is_staff` tinyint(1) NOT NULL,
  `is_active` tinyint(1) NOT NULL,
  `date_joined` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of auth_user
-- ----------------------------
INSERT INTO `auth_user` VALUES ('5', 'pbkdf2_sha256$36000$QLXvC5MYkd1d$W2Pm//LuH7aE8GxRCpYUVRH+kD45+kgEn4PXWWvVxso=', '2017-10-12 03:02:32', '0', 'admin', 'admin', 'admin', '', '0', '1', '2017-09-06 17:46:11');

-- ----------------------------
-- Table structure for auth_user_groups
-- ----------------------------
CREATE TABLE IF NOT EXISTS `auth_user_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `group_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth_user_groups_user_id_group_id_94350c0c_uniq` (`user_id`,`group_id`),
  KEY `auth_user_groups_group_id_97559544_fk_auth_group_id` (`group_id`),
  CONSTRAINT `auth_user_groups_group_id_97559544_fk_auth_group_id` FOREIGN KEY (`group_id`) REFERENCES `auth_group` (`id`),
  CONSTRAINT `auth_user_groups_user_id_6a12ed8b_fk_auth_user_id` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of auth_user_groups
-- ----------------------------

-- ----------------------------
-- Table structure for auth_user_user_permissions
-- ----------------------------
CREATE TABLE IF NOT EXISTS `auth_user_user_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `permission_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth_user_user_permissions_user_id_permission_id_14a6b632_uniq` (`user_id`,`permission_id`),
  KEY `auth_user_user_permi_permission_id_1fbb5f2c_fk_auth_perm` (`permission_id`),
  CONSTRAINT `auth_user_user_permissions_user_id_a95ead1b_fk_auth_user_id` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`),
  CONSTRAINT `auth_user_user_permi_permission_id_1fbb5f2c_fk_auth_perm` FOREIGN KEY (`permission_id`) REFERENCES `auth_permission` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of auth_user_user_permissions
-- ----------------------------

-- ----------------------------
-- Table structure for auths_userinfo
-- ----------------------------
CREATE TABLE IF NOT EXISTS `auths_userinfo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `phone` varchar(20) NOT NULL,
  `qq` varchar(20) NOT NULL,
  `weixin` varchar(20) NOT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`),
  CONSTRAINT `auths_userinfo_user_id_cff08a7a_fk_auth_user_id` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of auths_userinfo
-- ----------------------------

-- ----------------------------
-- Table structure for authtoken_token
-- ----------------------------
CREATE TABLE IF NOT EXISTS `authtoken_token` (
  `key` varchar(40) NOT NULL,
  `created` datetime NOT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`key`),
  UNIQUE KEY `user_id` (`user_id`),
  CONSTRAINT `authtoken_token_user_id_35299eff_fk_auth_user_id` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of authtoken_token
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_ci_attr
-- ----------------------------
CREATE TABLE IF NOT EXISTS `cmdb_ci_attr` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `display_name` varchar(255) DEFAULT NULL,
  `description` varchar(2000) DEFAULT NULL,
  `key` varchar(255) NOT NULL,
  `value` varchar(2000) NOT NULL,
  `mode` varchar(2) NOT NULL,
  `datatype` varchar(32) DEFAULT NULL,
  `intvalue` int(11) DEFAULT NULL,
  `stringvalue` varchar(2000) DEFAULT NULL,
  `datetimevalue` datetime DEFAULT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `ci_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `cmdb_ci_attr_ci_id_adcd572a_fk_cmdb_ci_object_id` (`ci_id`),
  CONSTRAINT `cmdb_ci_attr_ci_id_adcd572a_fk_cmdb_ci_object_id` FOREIGN KEY (`ci_id`) REFERENCES `cmdb_ci_object` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4101 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for cmdb_ci_attr_keylist
-- ----------------------------
CREATE TABLE IF NOT EXISTS `cmdb_ci_attr_keylist` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `key` varchar(2000) NOT NULL,
  `mode` varchar(1) NOT NULL,
  `create_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=69 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cmdb_ci_attr_keylist
-- ----------------------------
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('3', '主机名', 'SYSTEM.HOSTNAME', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('4', '默认mac地址', 'SYSTEM.DEFAULT.MAC', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('5', '默认IP（V4）', 'SYSTEM.DEFAULT.IP_V4', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('6', '默认IP（V6）', 'SYSTEM.DEFAULT.IP_V6', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('7', '默认网关', 'SYSTEM.DEFAULT.GATEWAY', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('8', 'CPU架构', 'SYSTEM.ARCHITECTURE', '2', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('9', 'OS版本', 'SYSTEM.OSVERSION', '2', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('10', 'OS家族', 'SYSTEM.OSFAMILY', '2', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('11', '网卡1 IP(V4)', 'SYSTEM.NIC.1.IP_V4', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('12', '网卡1 IP(V6)', 'SYSTEM.NIC.1.IP_V6', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('13', '网卡1 name', 'SYSTEM.NIC.1.NAME', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('14', '网卡1 广播域', 'SYSTEM.NIC.1.BROADCAST', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('15', '网卡1 mac地址', 'SYSTEM.NIC.1.MAC', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('16', '网卡1 子网掩码', 'SYSTEM.NIC.1.NETMASK', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('17', '网卡2 IP(V4)', 'SYSTEM.NIC.2.IP_V4', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('18', '网卡2 IP(V6)', 'SYSTEM.NIC.2.IP_V6', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('19', '网卡2 name', 'SYSTEM.NIC.2.NAME', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('20', '网卡2 广播域', 'SYSTEM.NIC.2.BROADCAST', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('21', '网卡2 mac地址', 'SYSTEM.NIC.2.MAC', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('22', '网卡2 子网掩码', 'SYSTEM.NIC.2.NETMASK', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('23', '内存（mb）可用', 'SYSTEM.MEM.MB.FREE', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('24', '内存（mb）已用', 'SYSTEM.MEM.MB.USED', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('25', '内存（mb）总共', 'SYSTEM.MEM.MB.TOTAL', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('26', '磁盘（gb）可用', 'SYSTEM.DISK.GB.FREE', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('27', '磁盘（gb）已用', 'SYSTEM.DISK.GB.USED', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('28', '磁盘（gb）总大小', 'SYSTEM.DISK.GB.SIZE', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('29', '磁盘sda（gb）可用', 'SYSTEM.DISK.SDA.GB.FREE', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('30', '磁盘sda（gb）已用', 'SYSTEM.DISK.SDA.GB.USED', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('31', '磁盘sda（gb）总大小', 'SYSTEM.DISK.SDA.GB.SIZE', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('32', 'bios版本', 'DEVICE.BIOS_VERSION', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('33', '设备品牌厂商', 'DEVICE.BRAND', '2', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('34', '设备序列号', 'DEVICE.SERIALNUMBER', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('35', '资产编号', 'DEVICE.ASSETNUMBER', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('36', '购买日期', 'DEVICE.DATE.PURCHASE', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('37', '上架日期', 'DEVICE.DATE.ONBOARD', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('38', '使用日期', 'DEVICE.DATE.MOVE2PRODUCTION', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('39', '过保日期', 'DEVICE.DATE.WARRANTY', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('40', '厂家联系方式', 'DEVICE.PROVIDERCONTRACTS', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('57', 'CONSUL地址', 'consul_addr', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('58', '报警状态', 'has_alarm', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('59', '英文名称', 'name', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('60', '名称', 'display_name', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('61', '上层资源', 'p_parent_name', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('62', '描述', 'description', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('63', '是否关联', 'is_link', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('65', '监控状态', 'mon_status', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('66', '运行状态', 'run_status', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('73', '交换机用户名', 'SWITCH_USER', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('74', '交换机密码', 'SWITCH_PASSWORD', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('75', '交换机IP', 'SWITCH_IP', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('76', '交换机community', 'SWITCH_COMMUNITY', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('77', '监控状态2', 'auto.status', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('84', 'SSH地址', 'SSH_HOST', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('85', 'SSH端口', 'SSH_PORT', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('86', 'SSH用户', 'SSH_USER', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('87', 'SSH密码', 'SSH_PASSWD', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('88', 'IPMI地址', 'IPMI_ADDR', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('89', 'IPMI用户', 'IPME_USER', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('90', 'IPMI密码', 'IPMI_PASSWD', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('91', '智能上报', 'link_name', '1', '0000-00-00 00:00:00');
INSERT INTO `cmdb_ci_attr_keylist` VALUES ('92', '资源类型', 'resource_type', '1', '0000-00-00 00:00:00');


-- ----------------------------
-- Table structure for cmdb_ci_attr_option
-- ----------------------------
CREATE TABLE IF NOT EXISTS `cmdb_ci_attr_option` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `attr_key` varchar(255) NOT NULL,
  `o_key` int(11) DEFAULT NULL,
  `o_value` varchar(255) NOT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cmdb_ci_attr_option
-- ----------------------------
INSERT INTO `cmdb_ci_attr_option` VALUES ('1', 'SYSTEM.ARCHITECTURE', '0', 'x86', '2017-11-13 08:53:28', '2017-11-13 08:53:28');
INSERT INTO `cmdb_ci_attr_option` VALUES ('2', 'SYSTEM.OSVERSION', '0', 'win10', '2017-11-13 08:53:43', '2017-11-13 08:53:43');
INSERT INTO `cmdb_ci_attr_option` VALUES ('3', 'SYSTEM.ARCHITECTURE', '0', 'power', '2017-11-13 10:50:39', '2017-11-13 10:50:39');

-- ----------------------------
-- Table structure for cmdb_ci_datacenter
-- ----------------------------
CREATE TABLE IF NOT EXISTS `cmdb_ci_datacenter` (
  `uuid` varchar(64) NOT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `display_name` varchar(255) NOT NULL,
  `description` varchar(2000) DEFAULT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for cmdb_ci_interface
-- ----------------------------
CREATE TABLE IF NOT EXISTS `cmdb_ci_interface` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `display_name` varchar(255) NOT NULL,
  `description` varchar(2000) DEFAULT NULL,
  `connection` int(11) DEFAULT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `ci_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `cmdb_ci_interface_ci_id_f0f7bec9_fk_cmdb_ci_object_id` (`ci_id`),
  CONSTRAINT `cmdb_ci_interface_ci_id_f0f7bec9_fk_cmdb_ci_object_id` FOREIGN KEY (`ci_id`) REFERENCES `cmdb_ci_object` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cmdb_ci_interface
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_ci_log
-- ----------------------------
CREATE TABLE IF NOT EXISTS `cmdb_ci_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `level` varchar(1) NOT NULL,
  `title` varchar(255) NOT NULL,
  `content` longtext NOT NULL,
  `submiter` varchar(255) NOT NULL,
  `ack` varchar(1) NOT NULL,
  `cid` int(11) DEFAULT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for cmdb_ci_log_follower
-- ----------------------------
CREATE TABLE IF NOT EXISTS `cmdb_ci_log_follower` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL,
  `content` longtext NOT NULL,
  `submiter` varchar(255) NOT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `logid_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `cmdb_ci_log_follower_logid_id_f71faa1c_fk_cmdb_ci_log_id` (`logid_id`),
  CONSTRAINT `cmdb_ci_log_follower_logid_id_f71faa1c_fk_cmdb_ci_log_id` FOREIGN KEY (`logid_id`) REFERENCES `cmdb_ci_log` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cmdb_ci_log_follower
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_ci_object
-- ----------------------------
CREATE TABLE IF NOT EXISTS `cmdb_ci_object` (
  `uuid` varchar(64) NOT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `display_name` varchar(255) NOT NULL,
  `description` varchar(4096) DEFAULT NULL,
  `status` varchar(2) NOT NULL,
  `mode` varchar(2) NOT NULL,
  `type_name` varchar(255) NOT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `type_id` int(11) NOT NULL,
  `zone_id` int(11) NOT NULL,
  `mon_status` varchar(2) NOT NULL,
  `p_parent_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `cmdb_ci_object_type_id_e625516b_fk_cmdb_ci_type_id` (`type_id`),
  KEY `cmdb_ci_object_zone_id_fea74747_fk_cmdb_ci_zone_id` (`zone_id`),
  KEY `cmdb_ci_object_p_parent_id_1e16095f_fk_cmdb_ci_object_id` (`p_parent_id`),
  CONSTRAINT `cmdb_ci_object_p_parent_id_1e16095f_fk_cmdb_ci_object_id` FOREIGN KEY (`p_parent_id`) REFERENCES `cmdb_ci_object` (`id`),
  CONSTRAINT `cmdb_ci_object_type_id_e625516b_fk_cmdb_ci_type_id` FOREIGN KEY (`type_id`) REFERENCES `cmdb_ci_type` (`id`),
  CONSTRAINT `cmdb_ci_object_zone_id_fea74747_fk_cmdb_ci_zone_id` FOREIGN KEY (`zone_id`) REFERENCES `cmdb_ci_zone` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for cmdb_ci_object_phy
-- ----------------------------
CREATE TABLE IF NOT EXISTS `cmdb_ci_object_phy` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `original_id_id` int(11) NOT NULL,
  `p_parent_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `original_id_id` (`original_id_id`),
  KEY `cmdb_ci_object_phy_p_parent_id_02ff6633_fk_cmdb_ci_object_phy_id` (`p_parent_id`),
  CONSTRAINT `cmdb_ci_object_phy_original_id_id_990c6b04_fk_cmdb_ci_object_id` FOREIGN KEY (`original_id_id`) REFERENCES `cmdb_ci_object` (`id`),
  CONSTRAINT `cmdb_ci_object_phy_p_parent_id_02ff6633_fk_cmdb_ci_object_phy_id` FOREIGN KEY (`p_parent_id`) REFERENCES `cmdb_ci_object_phy` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cmdb_ci_object_phy
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_ci_object_virtual
-- ----------------------------
CREATE TABLE  IF NOT EXISTS `cmdb_ci_object_virtual` (
  `id` varchar(64) NOT NULL,
  `name` varchar(255) NOT NULL,
  `display_name` varchar(255) NOT NULL,
  `description` varchar(2000) DEFAULT NULL,
  `is_root` tinyint(1) NOT NULL,
  `object_id` int(11) DEFAULT NULL,
  `sub_id` varchar(255) DEFAULT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `p_parent_id` varchar(64) DEFAULT NULL,
  `datacenter` varchar(255) DEFAULT NULL,
  `location_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `cmdb_ci_object_virtual_p_parent_id_b776160d_fk` (`p_parent_id`),
  KEY `cmdb_ci_object_virtual_location_id_a1a99918_fk_cmdb_ci_zone_id` (`location_id`),
  CONSTRAINT `cmdb_ci_object_virtual_location_id_a1a99918_fk_cmdb_ci_zone_id` FOREIGN KEY (`location_id`) REFERENCES `cmdb_ci_zone` (`id`),
  CONSTRAINT `cmdb_ci_object_virtual_p_parent_id_b776160d_fk` FOREIGN KEY (`p_parent_id`) REFERENCES `cmdb_ci_object_virtual` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cmdb_ci_object_virtual
-- ----------------------------

-- ----------------------------
-- Table structure for cmdb_ci_type
-- ----------------------------
CREATE TABLE  IF NOT EXISTS `cmdb_ci_type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `display_name` varchar(255) NOT NULL,
  `description` varchar(2000) DEFAULT NULL,
  `base_attr` varchar(2000) DEFAULT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cmdb_ci_type
-- ----------------------------
INSERT INTO `cmdb_ci_type` VALUES ('1', 'Rack', '机柜', 'Rack', null, '2017-07-19 06:51:45', '2017-07-19 06:51:45');
INSERT INTO `cmdb_ci_type` VALUES ('2', 'KVM', 'KVM虚拟机', 'Virtual server', null, '2017-07-19 06:51:45', '2017-07-19 06:51:45');
INSERT INTO `cmdb_ci_type` VALUES ('3', 'Switch', '交换机', 'Switch', null, '2017-07-19 06:51:45', '2017-07-19 06:51:45');
INSERT INTO `cmdb_ci_type` VALUES ('4', 'APP', '应用', '应用的分类', null, '2017-09-05 03:55:59', '2017-09-05 03:55:59');
INSERT INTO `cmdb_ci_type` VALUES ('5', 'Host', '物理机', 'Physical server', null, '2017-09-05 10:28:56', '2017-09-05 10:28:56');
INSERT INTO `cmdb_ci_type` VALUES ('6', 'VMware', 'VMware虚拟机', 'VMware', null, '2017-11-02 09:36:33', '2017-11-02 09:36:33');
INSERT INTO `cmdb_ci_type` VALUES ('7', 'Storage', '云硬盘', '云硬盘', null, '2017-12-07 04:26:49', '2017-12-07 04:26:49');

-- ----------------------------
-- Table structure for cmdb_ci_zone
-- ----------------------------
CREATE TABLE  IF NOT EXISTS `cmdb_ci_zone` (
  `uuid` varchar(64) NOT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `display_name` varchar(255) NOT NULL,
  `description` varchar(2000) DEFAULT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `dc_id` int(11) NOT NULL,
  `rank_length` int(11) NOT NULL,
  `rank_width` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `cmdb_ci_zone_dc_id_f36bcab0_fk_cmdb_ci_datacenter_id` (`dc_id`),
  CONSTRAINT `cmdb_ci_zone_dc_id_f36bcab0_fk_cmdb_ci_datacenter_id` FOREIGN KEY (`dc_id`) REFERENCES `cmdb_ci_datacenter` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for cmdb_table_view
-- ----------------------------
CREATE TABLE  IF NOT EXISTS `cmdb_table_view` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `default` tinyint(1) NOT NULL,
  `isedit` tinyint(1) NOT NULL,
  `keys` varchar(1024) NOT NULL,
  `dtype_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `cmdb_table_view_dtype_id_fd68c189_fk_cmdb_ci_type_id` (`dtype_id`),
  CONSTRAINT `cmdb_table_view_dtype_id_fd68c189_fk_cmdb_ci_type_id` FOREIGN KEY (`dtype_id`) REFERENCES `cmdb_ci_type` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8;



-- ----------------------------
-- Table structure for cmdb_type_link_attr_key
-- ----------------------------
CREATE TABLE  IF NOT EXISTS `cmdb_type_link_attr_key` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `attr_keylist_id` int(11) NOT NULL,
  `k_type_id` int(11) NOT NULL,
  `ischeck` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `cmdb_type_link_attr__attr_keylist_id_3fc1db9a_fk_cmdb_ci_a` (`attr_keylist_id`),
  KEY `cmdb_type_link_attr_key_k_type_id_20402ebe_fk_cmdb_ci_type_id` (`k_type_id`),
  CONSTRAINT `cmdb_type_link_attr_key_k_type_id_20402ebe_fk_cmdb_ci_type_id` FOREIGN KEY (`k_type_id`) REFERENCES `cmdb_ci_type` (`id`),
  CONSTRAINT `cmdb_type_link_attr__attr_keylist_id_3fc1db9a_fk_cmdb_ci_a` FOREIGN KEY (`attr_keylist_id`) REFERENCES `cmdb_ci_attr_keylist` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=208 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for cmdb_type_link_zone
-- ----------------------------
CREATE TABLE  IF NOT EXISTS `cmdb_type_link_zone` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `z_type_id` int(11) NOT NULL,
  `zone_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `cmdb_type_link_zone_z_type_id_b31bd46f_fk_cmdb_ci_type_id` (`z_type_id`),
  KEY `cmdb_type_link_zone_zone_id_d83f0e3d_fk_cmdb_ci_zone_id` (`zone_id`),
  CONSTRAINT `cmdb_type_link_zone_zone_id_d83f0e3d_fk_cmdb_ci_zone_id` FOREIGN KEY (`zone_id`) REFERENCES `cmdb_ci_zone` (`id`),
  CONSTRAINT `cmdb_type_link_zone_z_type_id_b31bd46f_fk_cmdb_ci_type_id` FOREIGN KEY (`z_type_id`) REFERENCES `cmdb_ci_type` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=97 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for django_admin_log
-- ----------------------------
CREATE TABLE  IF NOT EXISTS `django_admin_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `action_time` datetime NOT NULL,
  `object_id` longtext,
  `object_repr` varchar(200) NOT NULL,
  `action_flag` smallint(5) unsigned NOT NULL,
  `change_message` longtext NOT NULL,
  `content_type_id` int(11) DEFAULT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `django_admin_log_content_type_id_c4bce8eb_fk_django_co` (`content_type_id`),
  KEY `django_admin_log_user_id_c564eba6_fk_auth_user_id` (`user_id`),
  CONSTRAINT `django_admin_log_content_type_id_c4bce8eb_fk_django_co` FOREIGN KEY (`content_type_id`) REFERENCES `django_content_type` (`id`),
  CONSTRAINT `django_admin_log_user_id_c564eba6_fk_auth_user_id` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of django_admin_log
-- ----------------------------

-- ----------------------------
-- Table structure for django_content_type
-- ----------------------------
CREATE TABLE  IF NOT EXISTS `django_content_type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_label` varchar(100) NOT NULL,
  `model` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `django_content_type_app_label_model_76bd3d3b_uniq` (`app_label`,`model`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of django_content_type
-- ----------------------------
INSERT INTO `django_content_type` VALUES ('1', 'admin', 'logentry');
INSERT INTO `django_content_type` VALUES ('3', 'auth', 'group');
INSERT INTO `django_content_type` VALUES ('2', 'auth', 'permission');
INSERT INTO `django_content_type` VALUES ('4', 'auth', 'user');
INSERT INTO `django_content_type` VALUES ('7', 'authtoken', 'token');
INSERT INTO `django_content_type` VALUES ('14', 'cmdb', 'ci_attr');
INSERT INTO `django_content_type` VALUES ('17', 'cmdb', 'ci_attr_keylist');
INSERT INTO `django_content_type` VALUES ('8', 'cmdb', 'ci_attr_option');
INSERT INTO `django_content_type` VALUES ('19', 'cmdb', 'ci_datacenter');
INSERT INTO `django_content_type` VALUES ('12', 'cmdb', 'ci_interface');
INSERT INTO `django_content_type` VALUES ('11', 'cmdb', 'ci_log');
INSERT INTO `django_content_type` VALUES ('9', 'cmdb', 'ci_log_follower');
INSERT INTO `django_content_type` VALUES ('18', 'cmdb', 'ci_object');
INSERT INTO `django_content_type` VALUES ('15', 'cmdb', 'ci_object_phy');
INSERT INTO `django_content_type` VALUES ('16', 'cmdb', 'ci_object_virtual');
INSERT INTO `django_content_type` VALUES ('13', 'cmdb', 'ci_type');
INSERT INTO `django_content_type` VALUES ('10', 'cmdb', 'ci_zone');
INSERT INTO `django_content_type` VALUES ('22', 'cmdb', 'tableview');
INSERT INTO `django_content_type` VALUES ('23', 'cmdb', 'table_item');
INSERT INTO `django_content_type` VALUES ('24', 'cmdb', 'table_view');
INSERT INTO `django_content_type` VALUES ('21', 'cmdb', 'type_link_attr_key');
INSERT INTO `django_content_type` VALUES ('20', 'cmdb', 'type_link_zone');
INSERT INTO `django_content_type` VALUES ('5', 'contenttypes', 'contenttype');
INSERT INTO `django_content_type` VALUES ('6', 'sessions', 'session');

-- ----------------------------
-- Table structure for django_migrations
-- ----------------------------
CREATE TABLE  IF NOT EXISTS `django_migrations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `applied` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of django_migrations
-- ----------------------------
INSERT INTO `django_migrations` VALUES ('1', 'cmdb', '0001_initial', '2017-09-05 20:09:57');
INSERT INTO `django_migrations` VALUES ('2', 'contenttypes', '0001_initial', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('3', 'auth', '0001_initial', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('4', 'admin', '0001_initial', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('5', 'admin', '0002_logentry_remove_auto_add', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('6', 'contenttypes', '0002_remove_content_type_name', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('7', 'auth', '0002_alter_permission_name_max_length', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('8', 'auth', '0003_alter_user_email_max_length', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('9', 'auth', '0004_alter_user_username_opts', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('10', 'auth', '0005_alter_user_last_login_null', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('11', 'auth', '0006_require_contenttypes_0002', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('12', 'auth', '0007_alter_validators_add_error_messages', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('13', 'auth', '0008_alter_user_username_max_length', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('14', 'authtoken', '0001_initial', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('15', 'authtoken', '0002_auto_20160226_1747', '2017-09-05 20:10:12');
INSERT INTO `django_migrations` VALUES ('16', 'sessions', '0001_initial', '2017-09-05 20:10:13');
INSERT INTO `django_migrations` VALUES ('17', 'cmdb', '0002_auto_20170906_0418', '2017-09-05 20:18:09');
INSERT INTO `django_migrations` VALUES ('18', 'cmdb', '0003_auto_20170908_1901', '2017-09-08 11:01:35');
INSERT INTO `django_migrations` VALUES ('19', 'cmdb', '0004_auto_20170913_1929', '2017-09-13 11:30:18');
INSERT INTO `django_migrations` VALUES ('20', 'cmdb', '0002_auto_20170905_1840', '2017-10-19 03:05:22');
INSERT INTO `django_migrations` VALUES ('21', 'cmdb', '0002_type_link_attr_key_ischeck', '2017-11-22 06:56:41');
INSERT INTO `django_migrations` VALUES ('22', 'cmdb', '0003_auto_20171201_1554', '2017-12-01 07:54:15');
INSERT INTO `django_migrations` VALUES ('23', 'cmdb', '0004_auto_20171212_1817', '2017-12-12 10:18:05');
INSERT INTO `django_migrations` VALUES ('24', 'cmdb', '0005_auto_20171213_1127', '2017-12-13 03:27:29');
INSERT INTO `django_migrations` VALUES ('25', 'cmdb', '0006_table_view_isedit', '2017-12-13 08:51:24');
INSERT INTO `django_migrations` VALUES ('26', 'cmdb', '0002_auto_20171215_1733', '2017-12-15 09:33:31');

-- ----------------------------
-- Table structure for django_session
-- ----------------------------
CREATE TABLE  IF NOT EXISTS `django_session` (
  `session_key` varchar(40) NOT NULL,
  `session_data` longtext NOT NULL,
  `expire_date` datetime NOT NULL,
  PRIMARY KEY (`session_key`),
  KEY `django_session_expire_date_a5c62663` (`expire_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of django_session
-- ----------------------------
INSERT INTO `django_session` VALUES ('3ybs2eap5zpxrx2h419j0wrerb7igdfw', 'OWRkOTVkNDE3MWExODUyMTEwZDY5MTgxMGRkYmU5M2VhMTk0MDBlNjp7Il9hdXRoX3VzZXJfYmFja2VuZCI6ImRqYW5nby5jb250cmliLmF1dGguYmFja2VuZHMuTW9kZWxCYWNrZW5kIiwiX2F1dGhfdXNlcl9pZCI6IjUiLCJfYXV0aF91c2VyX2hhc2giOiIwZTgxYmU3ODc2OGU1OGNkNzZlYTc4YTdjYzYxNThlOGUyYjA4OGJiIn0=', '2017-09-28 16:13:58');
INSERT INTO `django_session` VALUES ('hkfsuxsxvqbo4z23oukh02civxz4ank7', 'NTI0MmRjZTkyOTJiN2JlMzhjYjgxY2M5ZjQ5NjYzOTE0ZDBjNjU5Zjp7Il9hdXRoX3VzZXJfaGFzaCI6IjBlODFiZTc4NzY4ZTU4Y2Q3NmVhNzhhN2NjNjE1OGU4ZTJiMDg4YmIiLCJfYXV0aF91c2VyX2lkIjoiNSIsIl9hdXRoX3VzZXJfYmFja2VuZCI6ImRqYW5nby5jb250cmliLmF1dGguYmFja2VuZHMuTW9kZWxCYWNrZW5kIn0=', '2017-10-26 03:02:32');

-- ----------------------------
-- Table structure for django_site
-- ----------------------------
CREATE TABLE IF NOT EXISTS `django_site` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `domain` varchar(100) NOT NULL,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `django_site_domain_a2e37b91_uniq` (`domain`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of django_site
-- ----------------------------
INSERT INTO `django_site` VALUES ('1', 'example.com', 'example.com');
