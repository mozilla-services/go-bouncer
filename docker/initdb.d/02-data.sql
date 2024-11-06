
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

LOCK TABLES `auth_group` WRITE;
/*!40000 ALTER TABLE `auth_group` DISABLE KEYS */;
INSERT INTO `auth_group` (`id`, `name`) VALUES (1,'API access');
INSERT INTO `auth_group` (`id`, `name`) VALUES (2,'Volunteer admins');
/*!40000 ALTER TABLE `auth_group` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `auth_group_permissions` WRITE;
/*!40000 ALTER TABLE `auth_group_permissions` DISABLE KEYS */;
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (1,1,40);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (2,1,41);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (3,1,42);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (4,1,43);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (5,1,44);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (6,1,45);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (7,1,46);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (8,1,47);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (9,1,48);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (10,1,49);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (19,2,31);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (11,2,32);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (12,2,33);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (13,2,34);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (14,2,35);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (15,2,36);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (16,2,37);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (17,2,38);
INSERT INTO `auth_group_permissions` (`id`, `group_id`, `permission_id`) VALUES (18,2,39);
/*!40000 ALTER TABLE `auth_group_permissions` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `auth_message` WRITE;
/*!40000 ALTER TABLE `auth_message` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_message` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `auth_permission` WRITE;
/*!40000 ALTER TABLE `auth_permission` DISABLE KEYS */;
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (1,'Can add Legacy User',1,'add_legacyuser');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (2,'Can change Legacy User',1,'change_legacyuser');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (3,'Can delete Legacy User',1,'delete_legacyuser');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (4,'Can add user profile',2,'add_userprofile');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (5,'Can change user profile',2,'change_userprofile');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (6,'Can delete user profile',2,'delete_userprofile');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (7,'Can add migration history',3,'add_migrationhistory');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (8,'Can change migration history',3,'change_migrationhistory');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (9,'Can delete migration history',3,'delete_migrationhistory');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (10,'Can add log entry',4,'add_logentry');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (11,'Can change log entry',4,'change_logentry');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (12,'Can delete log entry',4,'delete_logentry');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (13,'Can add permission',5,'add_permission');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (14,'Can change permission',5,'change_permission');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (15,'Can delete permission',5,'delete_permission');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (16,'Can add group',6,'add_group');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (17,'Can change group',6,'change_group');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (18,'Can delete group',6,'delete_group');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (19,'Can add user',7,'add_user');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (20,'Can change user',7,'change_user');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (21,'Can delete user',7,'delete_user');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (22,'Can add message',8,'add_message');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (23,'Can change message',8,'change_message');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (24,'Can delete message',8,'delete_message');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (25,'Can add content type',9,'add_contenttype');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (26,'Can change content type',9,'change_contenttype');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (27,'Can delete content type',9,'delete_contenttype');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (28,'Can add session',10,'add_session');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (29,'Can change session',10,'change_session');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (30,'Can delete session',10,'delete_session');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (31,'Can add country',11,'add_country');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (32,'Can change country',11,'change_country');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (33,'Can delete country',11,'delete_country');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (34,'Can add region',12,'add_region');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (35,'Can change region',12,'change_region');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (36,'Can delete region',12,'delete_region');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (37,'Can add mirror',13,'add_mirror');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (38,'Can change mirror',13,'change_mirror');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (39,'Can delete mirror',13,'delete_mirror');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (40,'Can add product',14,'add_product');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (41,'Can change product',14,'change_product');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (42,'Can delete product',14,'delete_product');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (43,'Can add product language',15,'add_product_language');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (44,'Can change product language',15,'change_product_language');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (45,'Can delete product language',15,'delete_product_language');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (46,'Can add location',16,'add_location');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (47,'Can change location',16,'change_location');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (48,'Can delete location',16,'delete_location');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (49,'Can view mirror uptake',16,'view_uptake');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (50,'Can add Operating System',17,'add_os');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (51,'Can change Operating System',17,'change_os');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (52,'Can delete Operating System',17,'delete_os');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (53,'Can add product alias',18,'add_productalias');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (54,'Can change product alias',18,'change_productalias');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (55,'Can delete product alias',18,'delete_productalias');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (56,'Can add product language',15,'add_productlanguage');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (57,'Can change product language',15,'change_productlanguage');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (58,'Can delete product language',15,'delete_productlanguage');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (59,'Can add location mirror map',19,'add_locationmirrormap');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (60,'Can change location mirror map',19,'change_locationmirrormap');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (61,'Can delete location mirror map',19,'delete_locationmirrormap');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (62,'Can add location mirror language exception',20,'add_locationmirrorlanguageexception');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (63,'Can change location mirror language exception',20,'change_locationmirrorlanguageexception');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (64,'Can delete location mirror language exception',20,'delete_locationmirrorlanguageexception');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (65,'Can add IP Block',21,'add_ipblock');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (66,'Can change IP Block',21,'change_ipblock');
INSERT INTO `auth_permission` (`id`, `name`, `content_type_id`, `codename`) VALUES (67,'Can delete IP Block',21,'delete_ipblock');
/*!40000 ALTER TABLE `auth_permission` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `auth_user` WRITE;
/*!40000 ALTER TABLE `auth_user` DISABLE KEYS */;
INSERT INTO `auth_user` (`id`, `username`, `first_name`, `last_name`, `email`, `password`, `is_staff`, `is_active`, `is_superuser`, `last_login`, `date_joined`) VALUES (1,'admin','','','admin@admin.com','pbkdf2_sha256$10000$1obyKBfQOSs6$TwvoLHHnE7uZprw9ZMZmviCVPKqCH1M+bFN2o6zNX4w=',1,1,1,'2015-07-23 16:07:01','2015-07-23 10:58:11');
/*!40000 ALTER TABLE `auth_user` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `auth_user_groups` WRITE;
/*!40000 ALTER TABLE `auth_user_groups` DISABLE KEYS */;
INSERT INTO `auth_user_groups` (`id`, `user_id`, `group_id`) VALUES (1,1,1);
/*!40000 ALTER TABLE `auth_user_groups` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `auth_user_user_permissions` WRITE;
/*!40000 ALTER TABLE `auth_user_user_permissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_user_user_permissions` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `django_admin_log` WRITE;
/*!40000 ALTER TABLE `django_admin_log` DISABLE KEYS */;
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (1,'2015-07-23 13:31:08',1,14,'1','Firefox',1,'');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (2,'2015-07-23 13:31:56',1,17,'1','win64',1,'');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (3,'2015-07-23 13:32:02',1,17,'2','osx',1,'');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (4,'2015-07-23 13:35:34',1,16,'1','/firefox/releases/40.0b5/win32/:lang/Firefox%20Setup%2040.0b5.exe',1,'');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (5,'2015-07-23 13:36:20',1,16,'2','/firefox/releases/40.0b5/mac/:lang/Firefox%2040.0b5.dmg',1,'');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (6,'2015-07-23 13:38:15',1,12,'1','All',1,'');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (7,'2015-07-23 13:38:30',1,13,'1','Mozilla Installer CDN',1,'');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (8,'2015-07-23 13:38:59',1,13,'2','Mozilla Installer CDN - SSL ',1,'');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (9,'2015-07-23 13:39:20',1,18,'1','ProductAlias object',1,'');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (10,'2015-07-23 13:41:49',1,11,'US','United States (US)',1,'');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (11,'2015-07-23 13:43:43',1,21,'1','0.0.0.0 -- 223.255.247.255',1,'');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (12,'2015-07-23 16:11:12',1,16,'2','/firefox/releases/39.0/mac/:lang/Firefox%2039.0.dmg',2,'Changed path.');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (13,'2015-07-23 16:11:28',1,16,'1','/firefox/releases/39.0/win32/:lang/Firefox%20Setup%2039.0.exe',2,'Changed path.');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (16,'2015-07-23 16:39:18',1,14,'2','Firefox-SSL',1,'');
INSERT INTO `django_admin_log` (`id`, `action_time`, `user_id`, `content_type_id`, `object_id`, `object_repr`, `action_flag`, `change_message`) VALUES (17,'2015-07-23 16:39:57',1,16,'3','/firefox/releases/39.0/win32/:lang/Firefox%20Setup%2039.0.exe',1,'');
/*!40000 ALTER TABLE `django_admin_log` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `django_content_type` WRITE;
/*!40000 ALTER TABLE `django_content_type` DISABLE KEYS */;
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (1,'Legacy User','users','legacyuser');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (2,'user profile','users','userprofile');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (3,'migration history','south','migrationhistory');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (4,'log entry','admin','logentry');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (5,'permission','auth','permission');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (6,'group','auth','group');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (7,'user','auth','user');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (8,'message','auth','message');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (9,'content type','contenttypes','contenttype');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (10,'session','sessions','session');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (11,'Country','geoip','country');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (12,'Region','geoip','region');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (13,'Mirror','mirror','mirror');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (14,'Product','mirror','product');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (15,'ProductLanguage','mirror','productlanguage');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (16,'Location','mirror','location');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (17,'Operating System','mirror','os');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (18,'product alias','mirror','productalias');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (19,'location mirror map','mirror','locationmirrormap');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (20,'location mirror language exception','mirror','locationmirrorlanguageexception');
INSERT INTO `django_content_type` (`id`, `name`, `app_label`, `model`) VALUES (21,'IP Block','geoip','ipblock');
/*!40000 ALTER TABLE `django_content_type` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `django_session` WRITE;
/*!40000 ALTER TABLE `django_session` DISABLE KEYS */;
INSERT INTO `django_session` (`session_key`, `session_data`, `expire_date`) VALUES ('810f108fd09ce4e4bb84684cb3581b41','OGJlMDI1NjU2YWM1YTA0NTAxY2M1YTU4NTkxYWIyYTY2MzA1ZTAxYzqAAn1xAShVEl9hdXRoX3Vz\nZXJfYmFja2VuZHECVSlkamFuZ28uY29udHJpYi5hdXRoLmJhY2tlbmRzLk1vZGVsQmFja2VuZHED\nVQ1fYXV0aF91c2VyX2lkcQSKAQF1Lg==\n','2015-08-06 16:07:01');
/*!40000 ALTER TABLE `django_session` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `geoip_country_to_region` WRITE;
/*!40000 ALTER TABLE `geoip_country_to_region` DISABLE KEYS */;
INSERT INTO `geoip_country_to_region` (`country_code`, `region_id`, `country_name`, `continent`) VALUES ('US',1,'United States','NA');
/*!40000 ALTER TABLE `geoip_country_to_region` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `geoip_ip_to_country` WRITE;
/*!40000 ALTER TABLE `geoip_ip_to_country` DISABLE KEYS */;
INSERT INTO `geoip_ip_to_country` (`id`, `ip_start`, `ip_end`, `country_code`) VALUES (1,0,3758094335,'US');
/*!40000 ALTER TABLE `geoip_ip_to_country` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `geoip_mirror_region_map` WRITE;
/*!40000 ALTER TABLE `geoip_mirror_region_map` DISABLE KEYS */;
INSERT INTO `geoip_mirror_region_map` (`id`, `mirror_id`, `region_id`) VALUES (1,1,1);
INSERT INTO `geoip_mirror_region_map` (`id`, `mirror_id`, `region_id`) VALUES (2,2,1);
/*!40000 ALTER TABLE `geoip_mirror_region_map` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `geoip_regions` WRITE;
/*!40000 ALTER TABLE `geoip_regions` DISABLE KEYS */;
INSERT INTO `geoip_regions` (`id`, `name`, `priority`, `throttle`, `fallback_id`, `prevent_global_fallback`) VALUES (1,'All',100,100,NULL,0);
/*!40000 ALTER TABLE `geoip_regions` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `mirror_aliases` WRITE;
/*!40000 ALTER TABLE `mirror_aliases` DISABLE KEYS */;
INSERT INTO `mirror_aliases` (`id`, `alias`, `related_product`) VALUES (1,'firefox-latest','Firefox');
INSERT INTO `mirror_aliases` (`id`, `alias`, `related_product`) VALUES (2,'firefox-sha1','Firefox-43.0.1-SSL');
INSERT INTO `mirror_aliases` (`id`, `alias`, `related_product`) VALUES (3,'firefox-beta-latest-ssl','Firefox-SSL');
INSERT INTO `mirror_aliases` (`id`, `alias`, `related_product`) VALUES (4,'firefox-devedition-latest-ssl','Devedition-128.0b1-SSL');
INSERT INTO `mirror_aliases` (`id`, `alias`, `related_product`) VALUES (5,'firefox-beta-latest','Firefox');
INSERT INTO `mirror_aliases` (`id`, `alias`, `related_product`) VALUES (6,'firefox-devedition-latest','Devedition-128.0b1');
INSERT INTO `mirror_aliases` (`id`, `alias`, `related_product`) VALUES (7,'firefox-latest-ssl','Firefox-SSL');
INSERT INTO `mirror_aliases` (`id`, `alias`, `related_product`) VALUES (8,'partner-firefox-release-unitedinternet-foo-latest','Firefox-partner-unitedinternet-foo');
INSERT INTO `mirror_aliases` (`id`, `alias`, `related_product`) VALUES (9,'firefox-beta-stub','Firefox-stub');
INSERT INTO `mirror_aliases` (`id`, `alias`, `related_product`) VALUES (10,'firefox-devedition-stub','Firefox-stub');
INSERT INTO `mirror_aliases` (`id`, `alias`, `related_product`) VALUES (11,'firefox-devedition-msi-latest-ssl','Firefox-beta-msi-latest-SSL');
INSERT INTO `mirror_aliases` (`id`, `alias`, `related_product`) VALUES (12,'thunderbird-latest-ssl','Thunderbird-SSL');
/*!40000 ALTER TABLE `mirror_aliases` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `mirror_lmm_lang_exceptions` WRITE;
/*!40000 ALTER TABLE `mirror_lmm_lang_exceptions` DISABLE KEYS */;
/*!40000 ALTER TABLE `mirror_lmm_lang_exceptions` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `mirror_location_mirror_map` WRITE;
/*!40000 ALTER TABLE `mirror_location_mirror_map` DISABLE KEYS */;
/*!40000 ALTER TABLE `mirror_location_mirror_map` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `mirror_locations` WRITE;
/*!40000 ALTER TABLE `mirror_locations` DISABLE KEYS */;
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/39.0/win64/:lang/Firefox%20Setup%2039.0.exe',1,1,1);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/39.0/mac/:lang/Firefox%2039.0.dmg',1,2,2);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/39.0/win32/:lang/Firefox%20Setup%2039.0.exe',1,3,3);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/39.0/win64/:lang/Firefox%20Setup%2039.0.exe',2,1,4);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/39.0/mac/:lang/Firefox%2039.0.dmg',2,2,5);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/39.0/win32/:lang/Firefox%20Setup%2039.0.exe',2,3,6);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/43.0.1/win64/:lang/Firefox%20Setup%2043.0.1.exe',3,1,7);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/43.0.1/mac/:lang/Firefox%2043.0.1.dmg',3,2,8);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/43.0.1/win32/:lang/Firefox%20Setup%2043.0.1.exe',3,3,9);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/latest-mozilla-central/firefox-128.0a1.:lang.win64.installer.exe',4,1,10);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/latest-mozilla-central/firefox-128.0a1.:lang.win32.installer.exe',4,3,11);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/latest-mozilla-central-l10n/firefox-128.0a1.:lang.win64.installer.exe',5,1,12);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/latest-mozilla-central-l10n/firefox-128.0a1.:lang.win32.installer.exe',5,3,13);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/latest-mozilla-central-l10n/firefox-128.0a1.:lang.win64.installer.exe',6,1,14);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/latest-mozilla-central-l10n/firefox-128.0a1.:lang.win32.installer.exe',6,3,15);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/latest-mozilla-central-l10n/firefox-128.0a1.:lang.win64.installer.exe',7,1,16);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/latest-mozilla-central-l10n/firefox-128.0a1.:lang.win32.installer.exe',7,3,17);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/2024/05/2024-05-06-09-48-55-mozilla-central-l10n/firefox-127.0a1.:lang.win64.installer.exe',8,1,18);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/2024/05/2024-05-06-09-48-55-mozilla-central-l10n/firefox-127.0a1.:lang.win32.installer.exe',8,3,19);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/2024/05/2024-05-06-09-48-55-mozilla-central-l10n/firefox-127.0a1.:lang.win64.installer.exe',9,1,20);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/2024/05/2024-05-06-09-48-55-mozilla-central-l10n/firefox-127.0a1.:lang.win32.installer.exe',9,3,21);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/127.0b9/win64/:lang/Firefox%20Setup%20127.0b9.exe',10,1,22);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/127.0b9/mac/:lang/Firefox%20Setup%20127.0b9.exe',10,2,23);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/127.0b9/win32/:lang/Firefox%20Setup%20127.0b9.exe',10,3,24);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/127.0b9/win64/:lang/Firefox%20Setup%20127.0b9.exe',11,1,25);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/127.0b9/mac/:lang/Firefox%20Setup%20127.0b9.exe',11,2,26);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/127.0b9/win32/:lang/Firefox%20Setup%20127.0b9.exe',11,3,27);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/devedition/releases/128.0b1/win64/:lang/Firefox%20Setup%20128.0b1.exe',12,1,28);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/devedition/releases/128.0b1/mac/:lang/Firefox%20Setup%20128.0b1.exe',12,2,29);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/devedition/releases/128.0b1/win32/:lang/Firefox%20Setup%20128.0b1.exe',12,3,30);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/devedition/releases/128.0b1/win64/:lang/Firefox%20Setup%20128.0b1.exe',13,1,31);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/devedition/releases/128.0b1/mac/:lang/Firefox%20Setup%20128.0b1.exe',13,2,32);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/devedition/releases/128.0b1/win32/:lang/Firefox%20Setup%20128.0b1.exe',13,3,33);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/devedition/releases/127.0b9/win64/:lang/Firefox%20Setup%20127.0b9.exe',14,1,34);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/devedition/releases/127.0b9/mac/:lang/Firefox%20Setup%20127.0b9.exe',14,2,35);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/devedition/releases/127.0b9/win32/:lang/Firefox%20Setup%20127.0b9.exe',14,3,36);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/devedition/releases/127.0b9/win64/:lang/Firefox%20Setup%20127.0b9.exe',15,1,37);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/devedition/releases/127.0b9/mac/:lang/Firefox%20Setup%20127.0b9.exe',15,2,38);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/devedition/releases/127.0b9/win32/:lang/Firefox%20Setup%20127.0b9.exe',15,3,39);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/127.0/win64/:lang/Firefox%20Setup%20127.0.exe',16,1,40);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/127.0/mac/:lang/Firefox%20Setup%20127.0.exe',16,2,41);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/127.0/win32/:lang/Firefox%20Setup%20127.0.exe',16,3,42);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/127.0/win64/:lang/Firefox%20Setup%20127.0.exe',17,1,43);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/127.0/mac/:lang/Firefox%20Setup%20127.0.exe',17,2,44);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/127.0/win32/:lang/Firefox%20Setup%20127.0.exe',17,3,45);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/partners/foo/bar/39.0/win64/:lang/Firefox%20Setup%2039.0.exe',18,1,46);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/partners/foo/bar/39.0/mac/:lang/Firefox%2039.0.dmg',18,2,47);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/partners/foo/bar/39.0/win32/:lang/Firefox%20Setup%2039.0.exe',18,3,48);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/partners/foo/bar/127.0/win64/:lang/Firefox%20Setup%20127.0.exe',19,1,49);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/partners/foo/bar/127.0/mac/:lang/Firefox%20127.0.dmg',19,2,50);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/partners/foo/bar/127.0/win32/:lang/Firefox%20Setup%20127.0.exe',19,3,51);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/115.16.1esr/win64/:lang/Firefox%20Setup%20115.16.1esr.exe',20,1,52);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/115.16.1esr/win32/:lang/Firefox%20Setup%20115.16.1esr.exe',20,3,53);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/131.0.3/win64/:lang/Firefox%20Setup%20131.0.3.msi',21,1,54);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/131.0.3/win32/:lang/Firefox%20Setup%20131.0.3.msi',21,3,55);

/* Those two locations have "win32" in their path on purpose */
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/131.0.3/win32/:lang/Firefox%20Installer.exe',22,1,56);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/131.0.3/win32/:lang/Firefox%20Installer.exe',22,3,57);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/latest-mozilla-central-l10n/Firefox%20Installer.en-US.exe',23,1,58);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/latest-mozilla-central-l10n/Firefox%20Installer.en-US.exe',23,3,59);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/128.3.1esr/win64/:lang/Firefox%20Setup%20128.3.1esr.exe',24,1,60);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/128.3.1esr/win32/:lang/Firefox%20Setup%20128.3.1esr.exe',24,3,61);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/132.0b9/win64/:lang/Firefox%20Setup%20132.0b9.msi',25,1,62);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/132.0b9/win32/:lang/Firefox%20Setup%20132.0b9.msi',25,3,63);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/latest-mozilla-central/firefox-133.0a1.en-US.win64.installer.msi',26,1,64);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/nightly/latest-mozilla-central/firefox-133.0a1.en-US.win32.installer.msi',26,3,65);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/128.3.1esr/win64/:lang/Firefox%20Setup%20128.3.1esr.msi',27,1,65);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/firefox/releases/128.3.1esr/win32/:lang/Firefox%20Setup%20128.3.1esr.msi',27,3,66);

INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/thunderbird/releases/131.0.1/win64/:lang/Thunderbird%20Setup%20131.0.1.exe',28,1,67);
INSERT INTO `mirror_locations` (`path`, `product_id`, `os_id`, `id`) VALUES ('/thunderbird/releases/131.0.1/win32/:lang/Thunderbird%20Setup%20131.0.1.exe',28,3,68);

/*!40000 ALTER TABLE `mirror_locations` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `mirror_mirrors` WRITE;
/*!40000 ALTER TABLE `mirror_mirrors` DISABLE KEYS */;
INSERT INTO `mirror_mirrors` (`count`, `rating`, `name`, `baseurl`, `active`, `id`) VALUES (0,100000,'Mozilla Installer CDN','http://download-installer.cdn.mozilla.net/pub',1,1);
INSERT INTO `mirror_mirrors` (`count`, `rating`, `name`, `baseurl`, `active`, `id`) VALUES (0,81000,'Mozilla Installer CDN - SSL ','https://download-installer.cdn.mozilla.net/pub',1,2);
/*!40000 ALTER TABLE `mirror_mirrors` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `mirror_mirrors_contacts` WRITE;
/*!40000 ALTER TABLE `mirror_mirrors_contacts` DISABLE KEYS */;
INSERT INTO `mirror_mirrors_contacts` (`id`, `mirror_id`, `user_id`) VALUES (1,1,1);
INSERT INTO `mirror_mirrors_contacts` (`id`, `mirror_id`, `user_id`) VALUES (2,2,1);
/*!40000 ALTER TABLE `mirror_mirrors_contacts` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `mirror_os` WRITE;
/*!40000 ALTER TABLE `mirror_os` DISABLE KEYS */;
INSERT INTO `mirror_os` (`priority`, `id`, `name`) VALUES (0,1,'win64');
INSERT INTO `mirror_os` (`priority`, `id`, `name`) VALUES (0,2,'osx');
INSERT INTO `mirror_os` (`priority`, `id`, `name`) VALUES (0,3,'win');
/*!40000 ALTER TABLE `mirror_os` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `mirror_product_langs` WRITE;
/*!40000 ALTER TABLE `mirror_product_langs` DISABLE KEYS */;
INSERT INTO `mirror_product_langs` (`language`, `product_id`, `id`) VALUES ('en-GB',1,1);
INSERT INTO `mirror_product_langs` (`language`, `product_id`, `id`) VALUES ('en-US',1,2);
INSERT INTO `mirror_product_langs` (`language`, `product_id`, `id`) VALUES ('en-US',2,3);
INSERT INTO `mirror_product_langs` (`language`, `product_id`, `id`) VALUES ('en-GB',2,4);
INSERT INTO `mirror_product_langs` (`language`, `product_id`, `id`) VALUES ('en-GB',3,5);
INSERT INTO `mirror_product_langs` (`language`, `product_id`, `id`) VALUES ('en-US',3,6);
/*!40000 ALTER TABLE `mirror_product_langs` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `mirror_products` WRITE;
/*!40000 ALTER TABLE `mirror_products` DISABLE KEYS */;
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox',1,1,1,1,0);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-SSL',1,1,1,2,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-43.0.1-SSL',1,1,1,3,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-nightly-latest-SSL',1,1,1,4,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-nightly-latest',1,1,1,5,0);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-nightly-latest-l10n-SSL',1,1,1,6,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-nightly-latest-l10n',1,1,1,7,0);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-nightly-pre2024-SSL',1,1,1,8,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-nightly-pre2024',1,1,1,9,0);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-127.0b9-SSL',1,1,1,10,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-127.0b9',1,1,1,11,0);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Devedition-128.0b1-SSL',1,1,1,12,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Devedition-128.0b1',1,1,1,13,0);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Devedition-127.0b9-SSL',1,1,1,14,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Devedition-127.0b9',1,1,1,15,0);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-127.0',1,1,1,16,0);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-127.0-SSL',1,1,1,17,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-partner-unitedinternet-foo',1,1,1,18,0);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-127.0-unitedinternet-foo',1,1,1,19,0);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-esr115-latest-SSL',1,1,1,20,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-msi-latest-SSL',1,1,1,21,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-stub',1,1,1,22,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-nightly-stub',1,1,1,23,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-esr-latest-SSL',1,1,1,24,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-beta-msi-latest-SSL',1,1,1,25,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-nightly-msi-latest-SSL',1,1,1,26,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Firefox-esr-msi-latest-SSL',1,1,1,27,1);
INSERT INTO `mirror_products` (`count`, `name`, `checknow`, `priority`, `active`, `id`, `ssl_only`) VALUES (0,'Thunderbird-SSL',1,1,1,28,1);
/*!40000 ALTER TABLE `mirror_products` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `sentry_log` WRITE;
/*!40000 ALTER TABLE `sentry_log` DISABLE KEYS */;
INSERT INTO `sentry_log` (`log_date`, `check_time`, `mirror_id`, `mirror_active`, `mirror_rating`, `reason`) VALUES ('2015-07-23 16:19:09','2015-07-23 23:19:09',1,'1',100000,'Checking mirror download-installer.cdn.mozilla.net ...\ndownload-installer.cdn.mozilla.net.	5	IN	CNAME	cs163.wpc.taucdn.net.\ncs163.wpc.taucdn.net.	5	IN	A	93.184.215.191\nUsing first seen IP: 93.184.215.191 for requests\nMaking base URL http://93.184.215.191/pub\n[2015-07-23 16:19:09 -0700] HEAD http://93.184.215.191/pub/firefox/releases/39.0/win32/zh-TW/Firefox%20Setup%2039.0.exe ... okay. CACHE=hit TOOK=0.006216\n[2015-07-23 16:19:09 -0700] HEAD http://93.184.215.191/pub/firefox/releases/39.0/mac/zh-TW/Firefox%2039.0.dmg ... okay. CACHE=hit TOOK=0.007065\nFinished. Elapsed time: 0.\n');
INSERT INTO `sentry_log` (`log_date`, `check_time`, `mirror_id`, `mirror_active`, `mirror_rating`, `reason`) VALUES ('2015-07-23 16:19:09','2015-07-23 23:19:28',2,'0',90000,'Checking mirror download-installer.cdn.mozilla.net ...\ndownload-installer.cdn.mozilla.net.	5	IN	CNAME	cs163.wpc.taucdn.net.\ncs163.wpc.taucdn.net.	5	IN	A	93.184.215.191\nUsing first seen IP: 93.184.215.191 for requests\nMaking base URL http://93.184.215.191/pub\n[2015-07-23 16:19:09 -0700] HEAD http://93.184.215.191/pub/firefox/releases/39.0/win32/zh-TW/Firefox%20Setup%2039.0.exe ... okay. CACHE=hit TOOK=0.006216\n[2015-07-23 16:19:09 -0700] HEAD http://93.184.215.191/pub/firefox/releases/39.0/mac/zh-TW/Firefox%2039.0.dmg ... okay. CACHE=hit TOOK=0.007065\nFinished. Elapsed time: 0.\nChecking mirror download-installer.cdn.mozilla.net ...\ndownload-installer.cdn.mozilla.net.	5	IN	CNAME	cs163.wpc.taucdn.net.\ncs163.wpc.taucdn.net.	5	IN	A	93.184.215.191\nUsing first seen IP: 93.184.215.191 for requests\nMaking base URL https://93.184.215.191/pub\nhttps://93.184.215.191/pub sent no response after 10 seconds!  Checking recent history...\n**** weight 90000 active 0 for https://93.184.215.191/pub\n**** https://93.184.215.191/pub Weight Drop Pattern matched, weight will be dropped 10%\n**** https://93.184.215.191/pub Weight change 90000 -> 81000\n');
/*!40000 ALTER TABLE `sentry_log` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `south_migrationhistory` WRITE;
/*!40000 ALTER TABLE `south_migrationhistory` DISABLE KEYS */;
INSERT INTO `south_migrationhistory` (`id`, `app_name`, `migration`, `applied`) VALUES (1,'mirror','0001_initial','2015-07-23 17:58:19');
INSERT INTO `south_migrationhistory` (`id`, `app_name`, `migration`, `applied`) VALUES (2,'mirror','0002_auto','2015-07-23 17:58:19');
INSERT INTO `south_migrationhistory` (`id`, `app_name`, `migration`, `applied`) VALUES (3,'mirror','0003_add_permissions_for_api','2015-07-23 17:58:19');
INSERT INTO `south_migrationhistory` (`id`, `app_name`, `migration`, `applied`) VALUES (4,'mirror','0004_auto__add_productalias__add_field_locationmirrormap_healthy__add_field','2015-07-23 17:58:20');
INSERT INTO `south_migrationhistory` (`id`, `app_name`, `migration`, `applied`) VALUES (5,'geoip','0001_initial','2015-07-23 17:58:20');
INSERT INTO `south_migrationhistory` (`id`, `app_name`, `migration`, `applied`) VALUES (6,'geoip','0002_auto__add_field_region_fallback','2015-07-23 17:58:20');
INSERT INTO `south_migrationhistory` (`id`, `app_name`, `migration`, `applied`) VALUES (7,'geoip','0003_auto__add_field_region_prevent_global_fallback','2015-07-23 17:58:20');
/*!40000 ALTER TABLE `south_migrationhistory` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `users_userprofile` WRITE;
/*!40000 ALTER TABLE `users_userprofile` DISABLE KEYS */;
INSERT INTO `users_userprofile` (`id`, `user_id`, `address`, `phone_number`, `ircnick`, `comments`) VALUES (1,1,'','','','');
/*!40000 ALTER TABLE `users_userprofile` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

