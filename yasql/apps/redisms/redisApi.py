# -*— coding: utf-8 -*-
# __author__ : pandonglin

import redis, time


REDIS_READ_CMDS = [  # 查redis询命令列表
    "help",
    "info",
    "exists",
    "keys",
    "type",
    "ttl",
    "scan",
    "get",
    "mget",
    "strlen",
    "hexists",
    "hget",
    "hlen",
    "hmget",
    "hvals",
    "hscan",
    "lindex",
    "llen",
    "lrange",
    "scard",
    "sismember",
    "smembers",
    "srandmember",
    "sscan",
    "zcard",
    "zcount",
    "zrange",
    "zrangebyscore",
    "zrank",
    "zrevrange",
    "zrevrangebyscore",
    "zrevrank",
    "zscore",
    "zscan",
]

REDIS_WRITE_CMDS = [  # 写redis命令列表,不包含删除动作的命令
    "set",
    "setex",
    "setnx",
    "hset",
    "hsetnx",
    "lpush",
    "lpushx",
    "rpush",
    "rpushx",
    "sadd",
    "zadd",
]

REDIS_CMDS = REDIS_READ_CMDS + REDIS_WRITE_CMDS


class RedisApi:
    def __init__(self, host, port, db=0, password=None):
        self.conn = self.get_conn(host, port, db, password)

    def get_conn(self, host, port, db, password):
        conn = redis.Redis(host, port, db, password=password, decode_responses=True)
        try:
            conn.ping()
        except Exception as err:
            raise Exception("can't connect redis")
        return conn

    def get_metrics(self, db):
        """监控指标"""
        result1 = self.conn.info()
        time.sleep(1)
        result2 = self.conn.info()
        data = {
            "version": result1.get("redis_version"),
            "run_time": "%s天" % result1.get("uptime_in_days"),
            "connected_clients": result1.get("connected_clients"),
            "blocked_clients": result1.get("blocked_clients"),
            "used_memory": result1.get("used_memory_human"),
            "total_memory": result1.get("total_system_memory_human"),
            "used_cpu_sys": float('%.2f' % (result2.get("used_cpu_sys") - result1.get("used_cpu_sys"))),
            "used_cpu_user": float('%.2f' % (result2.get("used_cpu_user") - result1.get("used_cpu_user"))),
            "hits": result1.get("keyspace_hits"),
            "misses": result1.get("keyspace_misses"),
            "ops_per_sec": result1.get("instantaneous_ops_per_sec"),
            "db_info": result1.get(db),
        }
        return data

    def read_help(self, args_list):
        if len(args_list) == 0:
            return REDIS_CMDS
        else:
            return "ERR wrong number of arguments for 'help' command"

    def read_info(self, args_list):
        if len(args_list) == 0:
            return self.conn.info()
        else:
            return "ERR wrong number of arguments for 'info' command"

    def read_exists(self, args_list):
        """检查给定 key 是否存在"""
        if len(args_list) == 1:
            try:
                r = self.conn.exists(args_list[0])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'exists' command"

    def read_keys(self, args_list):
        """查找所有符合给定模式 pattern 的 key"""
        if len(args_list) == 1 and args_list[0] != "*":
            try:
                r = self.conn.keys(args_list[0])
            except Exception as err:
                r = str(err)
            return r
        else:
            return " 'keys' command must have pattern，and pattern can't be '*' "

    def read_type(self, args_list):
        """返回 key 所储存的值的类型"""
        if len(args_list) == 1:
            return self.conn.type(args_list[0])
        else:
            return "ERR wrong number of arguments for 'type' command"

    def read_ttl(self, args_list):
        """返回 过期时间"""
        if len(args_list) == 1:
            return self.conn.ttl(args_list[0])
        else:
            return "ERR wrong number of arguments for 'ttl' command"

    def read_scan(self, args_list):
        """SCAN 命令及其相关的 SSCAN 命令、 HSCAN 命令和 ZSCAN 命令都用于增量地迭代（incrementally iterate）一集元素（a collection of elements"""
        try:
            if len(args_list) == 1:
                r = self.conn.scan(args_list[0])
            elif len(args_list) > 1 and len(args_list) % 2 != 0:
                args = args_list[1:]
                d = {}
                for i in range(len(args)):
                    if args[i].upper() == 'MATCH':
                        d['match'] = args[i+1]
                    elif args[i].upper() == 'COUNT':
                        d['count'] = args[i+1]
                r = self.conn.scan(args_list[0], **d)
            else:
                r = "ERR wrong arguments for 'scan' command"
        except Exception as err:
            r = str(err)
        return r

    def read_get(self, args_list):
        """返回 key 所关联的字符串值"""
        if len(args_list) == 1:
            try:
                r = self.conn.get(args_list[0])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'get' command"

    def read_mget(self, args_list):
        """返回所有(一个或多个)给定 key 的值"""
        if len(args_list) >= 1:
            try:
                r = self.conn.mget(args_list)
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'mget' command"

    def read_strlen(self, args_list):
        """返回 key 所储存的字符串值的长度"""
        if len(args_list) == 1:
            try:
                r = self.conn.strlen(args_list[0])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'strlen' command"

    def read_hexists(self, args_list):
        """查看哈希表 key 中，给定域 field 是否存在"""
        if len(args_list) == 2:
            try:
                r = self.conn.hexists(*args_list)
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'hexists' command"

    def read_hget(self, args_list):
        """返回哈希表 key 中给定域 field 的值"""
        if len(args_list) == 2:
            try:
                r = self.conn.hget(*args_list)
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'hget' command"

    def read_hlen(self, args_list):
        """返回哈希表 key 中域的数量"""
        if len(args_list) == 1:
            try:
                r = self.conn.hlen(args_list[0])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'hlen' command"

    def read_hmget(self, args_list):
        """返回哈希表 key 中，一个或多个给定域的值"""
        if len(args_list) >= 1:
            try:
                r = self.conn.hmget(args_list[0], args_list[1:])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'hmget' command"

    def read_hvals(self, args_list):
        """返回哈希表 key 中所有域的值"""
        if len(args_list) == 1:
            try:
                r = self.conn.hvals(args_list[0])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'hvals' command"

    def read_hscan(self, args_list):
        """用于增量地迭代"""
        try:
            if len(args_list) == 2:
                r = self.conn.hscan(args_list[0], args_list[1])
            elif len(args_list) > 2 and len(args_list) % 2 == 0:
                args = args_list[2:]
                d = {}
                for i in range(len(args)):
                    if args[i].upper() == 'MATCH':
                        d['match'] = args[i+1]
                    elif args[i].upper() == 'COUNT':
                        d['count'] = args[i+1]
                r = self.conn.hscan(args_list[0], args_list[1], **d)
            else:
                r = "ERR wrong arguments for 'hscan' command"
        except Exception as err:
            r = str(err)
        return r

    def read_lindex(self, args_list):
        """返回列表 key 中，下标为 index 的元素"""
        if len(args_list) == 2:
            try:
                r = self.conn.lindex(*args_list)
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'lindex' command"

    def read_llen(self, args_list):
        """返回列表 key 的长度"""
        if len(args_list) == 1:
            try:
                r = self.conn.llen(args_list[0])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'llen' command"

    def read_lrange(self, args_list):
        """返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定"""
        if len(args_list) == 3:
            try:
                r = self.conn.lrange(*args_list)
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'lrange' command"

    def read_scard(self, args_list):
        """返回集合 key 的基数(集合中元素的数量)"""
        if len(args_list) == 1:
            try:
                r = self.conn.scard(args_list[0])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'scard' command"

    def read_sismember(self, args_list):
        """判断 member 元素是否集合 key 的成员"""
        if len(args_list) == 2:
            try:
                r = self.conn.sismember(*args_list)
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'sismember' command"

    def read_smembers(self, args_list):
        """返回集合 key 中的所有成员"""
        if len(args_list) == 1:
            try:
                r = self.conn.smembers(args_list[0])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'smembers' command"

    def read_srandmember(self, args_list):
        """随机返回集合 key 中的成员"""
        if 1 <= len(args_list) <= 2:
            try:
                r = self.conn.srandmember(*args_list)
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'srandmember' command"

    def read_sscan(self, args_list):
        """用于增量地迭代"""
        try:
            if len(args_list) == 2:
                r = self.conn.sscan(args_list[0], args_list[1])
            elif len(args_list) > 2 and len(args_list) % 2 == 0:
                args = args_list[2:]
                d = {}
                for i in range(len(args)):
                    if args[i].upper() == 'MATCH':
                        d['match'] = args[i + 1]
                    elif args[i].upper() == 'COUNT':
                        d['count'] = args[i + 1]
                r = self.conn.sscan(args_list[0], args_list[1], **d)
            else:
                r = "ERR wrong arguments for 'sscan' command"
        except Exception as err:
            r = str(err)
        return r

    def read_zcard(self, args_list):
        """返回有序集 key 的基数"""
        if len(args_list) == 1:
            try:
                r = self.conn.zcard(args_list[0])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'zcard' command"

    def read_zcount(self, args_list):
        """返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量"""
        if len(args_list) == 3:
            try:
                r = self.conn.zcount(*args_list)
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'zcount' command"

    def read_zrange(self, args_list):  # TODO 待优化
        """返回有序集 key 中，指定区间内的成员"""
        try:
            if len(args_list) == 3:
                r = self.conn.zrange(*args_list)
            elif len(args_list) == 4 and args_list[3].upper() == "WITHSCORES":
                r = self.conn.zrange(args_list[0], args_list[1], args_list[2], withscores=True)
            else:
                r = "ERR wrong arguments for 'zrange' command"
        except Exception as err:
            r = str(err)
        return r

    def read_zrangebyscore(self, args_list):
        """ 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列 """
        try:
            if len(args_list) == 3:
                r = self.conn.zrangebyscore(*args_list)
            elif len(args_list) == 4 and args_list[3].upper() == "WITHSCORES":
                r = self.conn.zrangebyscore(args_list[0], args_list[1], args_list[2], withscores=True)
            elif len(args_list) == 7 and args_list[3].upper() == "WITHSCORES" and args_list[4] == "LIMIT":
                r = self.conn.zrangebyscore(args_list[0], args_list[1], args_list[2], args_list[5], args_list[6], withscores=True)
            else:
                r = "ERR wrong arguments for 'zrangebyscore' command"
        except Exception as err:
            r = str(err)
        return r

    def read_zrank(self, args_list):
        """返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列"""
        if len(args_list) == 2:
            try:
                r = self.conn.zrank(*args_list)
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'zrank' command"

    def read_zrevrange(self, args_list):
        """返回有序集 key 中，指定区间内的成员"""
        try:
            if len(args_list) == 3:
                r = self.conn.zrevrange(*args_list)
            elif len(args_list) == 4 and args_list[3].upper() == "WITHSCORES":
                r = self.conn.zrevrange(args_list[0], args_list[1], args_list[2], withscores=True)
            else:
                r = "ERR wrong arguments for 'zrevrange' command"
        except Exception as err:
            r = str(err)
        return r

    def read_zrevrangebyscore(self, args_list):
        """返回有序集 key 中， score 值介于 max 和 min 之间(默认包括等于 max 或 min )的所有的成员。有序集成员按 score 值递减(从大到小)的次序排列"""
        try:
            if len(args_list) == 3:
                r = self.conn.zrevrangebyscore(*args_list)
            elif len(args_list) == 4 and args_list[3].upper() == "WITHSCORES":
                r = self.conn.zrevrangebyscore(args_list[0], args_list[1], args_list[2], withscores=True)
            elif len(args_list) == 7 and args_list[3].upper() == "WITHSCORES" and args_list[4] == "LIMIT":
                r = self.conn.zrevrangebyscore(args_list[0], args_list[1], args_list[2], args_list[5], args_list[6], withscores=True)
            else:
                r = "ERR wrong arguments for 'zrevrangebyscore' command"
        except Exception as err:
            r = str(err)
        return r

    def read_zrevrank(self, args_list):
        """返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序"""
        if len(args_list) == 2:
            try:
                r = self.conn.zrevrank(*args_list)
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'zrevrank' command"

    def read_zscore(self, args_list):
        """返回有序集 key 中，成员 member 的 score 值"""
        if len(args_list) == 2:
            try:
                r = self.conn.zscore(*args_list)
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'zscore' command"

    def read_zscan(self, args_list):
        """用于增量地迭代"""
        try:
            if len(args_list) == 2:
                r = self.conn.zscan(args_list[0], args_list[1])
            elif len(args_list) > 2 and len(args_list) % 2 == 0:
                args = args_list[2:]
                d = {}
                for i in range(len(args)):
                    if args[i].upper() == 'MATCH':
                        d['match'] = args[i + 1]
                    elif args[i].upper() == 'COUNT':
                        d['count'] = args[i + 1]
                r = self.conn.zscan(args_list[0], args_list[1], **d)
            else:
                r = "ERR wrong arguments for 'zscan' command"
        except Exception as err:
            r = str(err)
        return r

    def write_set(self, args_list):
        """将字符串值 value 关联到 key"""
        try:
            if len(args_list) == 2:
                r = self.conn.set(args_list[0], args_list[1])
            elif len(args_list) == 3:
                d = {}
                if args_list[2].upper() == 'NX':
                    d['nx'] = True
                elif args_list[2].upper() == 'XX':
                    d['xx'] = True
                else:
                    return "ERR syntax error"
                r = self.conn.set(args_list[0], args_list[1], **d)
            elif len(args_list) > 3:
                args = args_list[2:]
                d = {}
                for i in range(len(args)):
                    if args[i].upper() == 'EX':
                        d['ex'] = args[i+1]
                    elif args[i].upper() == 'PX':
                        d['px'] = args[i+1]
                    elif args[i].upper() == 'NX':
                        d['nx'] = True
                    elif args[i].upper() == 'XX':
                        d['xx'] = True
                    else:
                        return "ERR syntax error"
                r = self.conn.set(args_list[0], args_list[1], **d)
            else:
                return "ERR wrong number of arguments for 'set' command"
        except Exception as err:
            r = str(err)
        return r

    def write_setex(self, args_list):
        """将值 value 关联到 key ，并将 key 的生存时间设为 seconds (以秒为单位)"""
        if len(args_list) == 3:
            try:
                r = self.conn.setex(args_list[0], args_list[2], args_list[1])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'setex' command"

    def write_setnx(self, args_list):
        """将 key 的值设为 value ，当且仅当 key 不存在"""
        if len(args_list) == 2:
            try:
                r = self.conn.setnx(args_list[0], args_list[1])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'setnx' command"

    def write_hset(self, args_list):
        """将哈希表 key 中的域 field 的值设为 value"""
        if len(args_list) == 3:
            try:
                r = self.conn.hset(args_list[0], args_list[1], args_list[2])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'hset' command"

    def write_hsetnx(self, args_list):
        if len(args_list) == 3:
            try:
                r = self.conn.hsetnx(args_list[0], args_list[1], args_list[2])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'hsetnx' command"

    def write_lpush(self, args_list):
        """将一个或多个值 value 插入到列表 key 的表头"""
        if len(args_list) >= 2:
            try:
                r = self.conn.lpush(args_list[0], *args_list[1:])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'lpush' command"

    def write_lpushx(self, args_list):
        """将值 value 插入到列表 key 的表头，当且仅当 key 存在并且是一个列表"""
        if len(args_list) == 2:
            try:
                r = self.conn.lpushx(args_list[0], args_list[1])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'lpushx' command"

    def write_rpush(self, args_list):
        """将一个或多个值 value 插入到列表 key 的表尾(最右边)"""
        if len(args_list) >= 2:
            try:
                r = self.conn.rpush(args_list[0], *args_list[1:])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'rpush' command"

    def write_rpushx(self, args_list):
        """将值 value 插入到列表 key 的表尾，当且仅当 key 存在并且是一个列表"""
        if len(args_list) == 2:
            try:
                r = self.conn.rpushx(args_list[0], args_list[1])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'rpushx' command"

    def write_sadd(self, args_list):
        """将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略"""
        if len(args_list) >= 2:
            try:
                r = self.conn.sadd(args_list[0], *args_list[1:])
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'sadd' command"

    def write_zadd(self, args_list):
        """将一个或多个 member 元素及其 score 值加入到有序集 key 当中"""
        if len(args_list) >= 2:
            try:
                r = self.conn.zadd(args_list[0], *reversed(args_list[1:]))
            except Exception as err:
                r = str(err)
            return r
        else:
            return "ERR wrong number of arguments for 'zadd' command"
