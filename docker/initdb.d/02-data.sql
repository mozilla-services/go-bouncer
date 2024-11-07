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

/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
