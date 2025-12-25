package model

const MysqlDSL = `CREATE TABLE IF NOT EXISTS %s (
 id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
 sink_id varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'sink_id',
 store_id varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '存储id，全局唯一 store_key+sink_type',
 hostname varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'hostname',
 crawl_url SMALLINT COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'crawl_url',
 crawl_source SMALLINT COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'crawl_source',
 crawl_type varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'crawl_type',
 data json NOT NULL COMMENT '客户端上报数据',
 metadata json DEFAULT NULL COMMENT '元数据',
 crawl_time datetime NOT NULL COMMENT '爬取时间',
 created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 PRIMARY KEY (id),
 UNIQUE KEY uniq_sink_store (sink_id, store_id),
 KEY idx_sink_id (sink_id),                      
  KEY idx_store_id (store_id)                   
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='爬虫记录表'
`

const PostgresDSL = `
CREATE TABLE IF NOT EXISTS %s (
    id BIGSERIAL PRIMARY KEY,
    sink_id TEXT NOT NULL,
    store_id TEXT NOT NULL,
    hostname TEXT NOT NULL,
    crawl_url TEXT NOT NULL,
    crawl_source SMALLINT NOT NULL,
    crawl_type SMALLINT NOT NULL,
    data JSONB NOT NULL,
    metadata JSONB,
    crawl_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uniq_sink_store UNIQUE (sink_id, store_id)
);

-- 创建额外的索引
CREATE INDEX IF NOT EXISTS idx_sink_id ON %s (sink_id);
CREATE INDEX IF NOT EXISTS idx_store_id ON %s (store_id);
`
