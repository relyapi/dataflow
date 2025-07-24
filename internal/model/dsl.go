package model

const MysqlDSL = `CREATE TABLE IF NOT EXISTS %s (
 id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
 sink_id varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'sink_id',
 store_id varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '存储id，全局唯一 store_key+sink_type',
 hostname varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'hostname',
 request_url varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'request_url',
 sink_type varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'sink_type',
 data json NOT NULL COMMENT '客户端上报数据',
 metadata json DEFAULT NULL COMMENT '元数据',
 crawl_time datetime NOT NULL COMMENT '爬取时间',
 created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
 PRIMARY KEY (id),
 UNIQUE KEY uniq_sink_store (sink_id, store_id),
 KEY idx_sink_id (sink_id),                      
  KEY idx_store_id (store_id)                   
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='爬虫记录表'
`
