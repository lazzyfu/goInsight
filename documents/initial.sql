INSERT INTO opsql.auditsql_role_permissions (id, permission_name, permission_desc, created_at, updated_at) VALUES (1, 'can_audit', '工单审核权限', '2019-03-18 03:24:41.577191', '2019-03-18 03:36:28.005257');
INSERT INTO opsql.auditsql_role_permissions (id, permission_name, permission_desc, created_at, updated_at) VALUES (3, 'can_commit', '工单提交权限', '2019-03-18 03:25:22.001932', '2019-03-18 03:26:01.509941');
INSERT INTO opsql.auditsql_role_permissions (id, permission_name, permission_desc, created_at, updated_at) VALUES (5, 'can_execute', '工单执行权限', '2019-03-18 03:26:31.232740', '2019-03-18 03:26:31.232770');


INSERT INTO opsql.auditsql_sys_environment (envi_id, envi_name, created_at, updated_at) VALUES (1, '测试环境', '2019-03-16 07:37:00.027384', '2019-03-16 07:37:00.027417');
INSERT INTO opsql.auditsql_sys_environment (envi_id, envi_name, created_at, updated_at) VALUES (2, '生产环境', '2019-03-16 07:37:05.062588', '2019-03-16 07:37:05.062618');


INSERT INTO opsql.auditsql_user_roles (rid, role_name, created_at, updated_at) VALUES (1, '开发', '2019-03-16 04:30:53.149371', '2019-03-21 02:13:51.369194');
INSERT INTO opsql.auditsql_user_roles (rid, role_name, created_at, updated_at) VALUES (2, 'DBA', '2019-03-18 03:36:11.352674', '2019-03-22 08:16:08.808813');
INSERT INTO opsql.auditsql_user_roles (rid, role_name, created_at, updated_at) VALUES (3, '测试', '2019-03-21 03:53:08.641261', '2019-03-28 03:04:00.287465');

INSERT INTO opsql.auditsql_useraccounts (password, last_login, is_superuser, username, first_name, last_name, email, is_staff, date_joined, uid, is_active, displayname, mobile, avatar_file) VALUES ('pbkdf2_sha256$120000$2ChYRZRJCJw1$4ADTQndkt32w4E7u1AtK/0mj1ZPBCDK81th9iWr9iNY=', '2019-04-09 02:22:34.863953', 1, 'admin', '', '', 'admin@example.com', 1, '2019-04-09 02:21:29.227011', 1, 1, '', null, 'img/avatar1.png');

INSERT INTO opsql.django_celery_beat_intervalschedule (id, every, period) VALUES (1, 1, 'minutes');
INSERT INTO opsql.django_celery_beat_periodictasks (ident, last_update) VALUES (1, '2019-04-09 04:51:00.367625');
INSERT INTO opsql.django_celery_beat_periodictask (id, name, task, args, kwargs, queue, exchange, routing_key, expires, enabled, last_run_at, total_run_count, date_changed, description, crontab_id, interval_id, solar_id, one_off, start_time, priority) VALUES (2, '同步远程库表结构元数据', 'query.tasks.periodic_sync_remote_schemameta_to_local', '[]', '{}', null, null, null, null, 1, null, 0, '2019-04-09 04:50:30.330153', '', null, 1, null, 0, null, null);
INSERT INTO opsql.django_celery_beat_periodictask (id, name, task, args, kwargs, queue, exchange, routing_key, expires, enabled, last_run_at, total_run_count, date_changed, description, crontab_id, interval_id, solar_id, one_off, start_time, priority) VALUES (3, '同步远程库元数据', 'orders.tasks.periodic_sync_schemas', '[]', '{}', null, null, null, null, 1, null, 0, '2019-04-09 04:51:00.369343', '', null, 1, null, 0, null, null);

INSERT INTO opsql.auditsql_useraccounts (password, last_login, is_superuser, username, first_name, last_name, email, is_staff, date_joined, uid, is_active, displayname, mobile, avatar_file) VALUES ('pbkdf2_sha256$120000$2ChYRZRJCJw1$4ADTQndkt32w4E7u1AtK/0mj1ZPBCDK81th9iWr9iNY=', '2019-04-09 04:49:02.610057', 1, 'admin', '', '', 'admin@example.com', 1, '2019-04-09 02:21:29.227011', 1, 1, '', null, 'img/avatar1.png');