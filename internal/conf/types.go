package conf

const MysqlDSL = `CREATE TABLE IF NOT EXISTS %s (
 	id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
 sink_id varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'sink_id',
 url_md5 varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'md5 url ',
 url text COLLATE utf8mb4_unicode_ci COMMENT '唯一url',
 data json NOT NULL COMMENT '差异字段（JSON格式）',
 crawl_time datetime NOT NULL COMMENT '爬取时间',
 created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
 PRIMARY KEY (id),
 UNIQUE KEY uniq_url_md5 (url_md5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='爬虫记录表'
`
