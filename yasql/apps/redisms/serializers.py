# -*- coding:utf-8 -*-
# by pandonglin
import shlex
from datetime import datetime
from rest_framework import serializers
from redisms import models, redisApi


class RedisGroupSerializers(serializers.ModelSerializer):
    class Meta:
        model = models.RedisGroup
        fields = "__all__"

    def to_representation(self, instance):
        ret = {
            "key": instance.id,
            "title": instance.display_full_name,
            "icon": "codepen",
            "children": [{"key": x.id, "title": x.full_dsn} for x in instance.rg.all()]
        }
        return ret


class RedisSerializer(serializers.Serializer):
    redis_id = serializers.IntegerField(required=True, error_messages={"required": "redis实例未选中"})
    redis_cmd = serializers.CharField(required=True, error_messages={"required": "没有可执行的命令"})
    redis_db = serializers.CharField(required=True, error_messages={"required": "没有可执行的实例"})

    def format_data(self, args):
        """格式化数据"""
        if isinstance(args, bytearray):
            return "bytearray data"
        elif isinstance(args, bytes):
            return args.decode("utf8", errors="ignore")
        elif isinstance(args, set):
            return list(args)
        elif isinstance(args, tuple):
            return list(args)
        else:
            return args

    def format_cmd(self, str):
        """格式化命令"""
        try:
            lexer = shlex.split(str)
            return lexer
        except:   # 不能解析的字符串，返回空
            pass

    def query(self, redis_api, cmd, args):
        f_cmd = "read_" + cmd
        if hasattr(redis_api, f_cmd):
            redis_func = getattr(redis_api, f_cmd)
            data = self.format_data(redis_func(args))
        else:
            data = "ERR unknown command '%s'" % cmd
        return data

    def execute(self, redis_api, cmd, args):
        f_cmd = "write_" + cmd
        if hasattr(redis_api, f_cmd):
            redis_func = getattr(redis_api, f_cmd)
            data = self.format_data(redis_func(args))
        else:
            data = "ERR unknown command '%s'" % cmd
        return data

    def handle(self, redis_cmd, redis_id, redis_db):
        cmd_args_list = self.format_cmd(redis_cmd)

        if cmd_args_list is None:  # 命令行解析失败
            raise serializers.ValidationError({"error": "Invalid argument(s)"})
        elif cmd_args_list[0].lower() not in redisApi.REDIS_CMDS:  # 不在可执行列表中
            raise serializers.ValidationError({"error": "禁止执行此命令"})
        else:
            try:
                redis_ins = models.RedisConfig.objects.get(id=redis_id)
                redis_api = redisApi.RedisApi(redis_ins.host, redis_ins.port, redis_db, redis_ins.decrypt_password)
                perm_list = redis_ins.group.rg_group.all()
            except Exception as err:
                raise serializers.ValidationError({"error": str(err)})
            else:
                # 判断读写权限
                if cmd_args_list[0].lower() in redisApi.REDIS_READ_CMDS and any([x.permission_code.startswith("read_") for x in perm_list]):
                    result = self.query(redis_api, cmd_args_list[0].lower(), cmd_args_list[1:])
                elif cmd_args_list[0].lower() in redisApi.REDIS_WRITE_CMDS and any([x.permission_code == "read_write" for x in perm_list]):
                    result = self.execute(redis_api, cmd_args_list[0].lower(), cmd_args_list[1:])
                else:
                    raise serializers.ValidationError({"error": "无权限执行"}, code='error')
                return result

    def validate(self, attrs):
        result = []
        cmd_list = attrs["redis_cmd"].split("\n")
        for cmd in cmd_list:
            result.append({"redis_cmd": cmd, "result": self.handle(cmd, attrs["redis_id"], attrs["redis_db"])})
        attrs["data"] = result
        return attrs

    def create(self, validated_data):
        # 记录日志
        request = self.context["request"]
        redis_ins = models.RedisConfig.objects.get(id=validated_data["redis_id"])
        models.RedisLog.objects.create(user=request.user.username, host_dsn=redis_ins.full_dsn, exec_cmd=validated_data["redis_cmd"])

        data = {
            "redis_name": redis_ins.full_dsn,
            "data": validated_data["data"]
        }
        return data