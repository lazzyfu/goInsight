SELECT id,mobile,TEXT,created_at        FROM sms_log WHERE
created_at >= '2018-07-17 00:00:00' AND
created_at <= '2018-07-17 23:59:59' AND TEXT LIKE '%客户%';

SELECT id,mobile,TEXT,created_at FROM     sms_log WHERE
created_at >= '2018-07-17 00:00:00' AND created_at <= '2018-07-18 23:59:59' AND TEXT LIKE '%客户%';


SELECT id,mobile,TEXT,created_at FROM sms_log WHERE  created_at >= '2018-07-17 00:00:00' AND created_at <= '2018-07-25 23:59:59' AND TEXT LIKE '%客户%';




