CREATE TABLE route(
                     `id` varchar(256) primary key,
                     `host` VARCHAR(256),
                     `port` int(64),
                     `path` VARCHAR(256),
                     `http_method` VARCHAR(256),
                     `response_header_template` text,
                     `response_body_template` text,
                     `http_template_id` VARCHAR(256)
);

INSERT INTO `mock`.`route` (`id`,`host`,`port`, `path`, `http_method`, `response_header_template`, `response_body_template`,`http_template_id`) VALUES (uuid(),'localhost:9999', 9999 , '/api/v1/test/mock', 'GET', NULL, '{
  "code": 200,
  "data": true,
  "message": "ok"
}',uuid());