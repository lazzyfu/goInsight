# -*— coding: utf-8 -*-
# __author__ : pandonglin
from collections import defaultdict
import shlex


class StatCounter(object):
    def __init__(self, prefix_delim=':'):
        self.line_count = 0
        self.skipped_lines = 0
        self.commands = defaultdict(int)
        self.keys = defaultdict(int)
        self.prefixes = defaultdict(int)
        self.times = []
        self._cached_sorts = {}
        self.start_ts = None
        self.last_ts = None
        self.prefix_delim = prefix_delim

    def _record_duration(self, entry):
        """处理执行时间"""
        ts = float(entry['time']) * 1000 * 1000    # microseconds
        if not self.start_ts:  # 收集的第1个命令仍掉
            self.start_ts = ts
            self.last_ts = ts
        else:
            duration = ts - self.last_ts
            self.times.append((duration, entry["command"]))
            self.last_ts = ts

    def _record_command(self, command):
        """统计执行命令"""
        self.commands[command] += 1

    def _record_key(self, key):
        """统计key"""
        self.keys[key] += 1
        parts = key.split(self.prefix_delim) # 按照前缀统计
        if len(parts) > 1:
            self.prefixes[parts[0]] += 1

    def _reformat_entry(self, entry):
        max_args_to_show = 5
        lexer = shlex.split(entry)
        output = '%s' % lexer[0]
        if len(lexer) > 1:
            output += ' %s' % lexer[1]
        if len(lexer) > 2:
            ellipses = ' ...' if len(lexer) > max_args_to_show else ''
            output += ' %s%s' % (' '.join(lexer[2:max_args_to_show]), ellipses)
        return output

    def _get_or_sort_list(self, ls):
        key = id(ls)
        if not key in self._cached_sorts:
            sorted_items = sorted(ls)
            self._cached_sorts[key] = sorted_items
        return self._cached_sorts[key]

    def _time_stats(self, times):
        sorted_times = self._get_or_sort_list(times)
        num_times = len(sorted_times)
        percent_50 = sorted_times[int(num_times / 2)][0]
        percent_75 = sorted_times[int(num_times * .75)][0]
        percent_90 = sorted_times[int(num_times * .90)][0]
        percent_99 = sorted_times[int(num_times * .99)][0]
        return [("Median", percent_50),
                ("75%", percent_75),
                ("90%", percent_90),
                ("99%", percent_99)]

    def _heaviest_commands(self, times):
        times_by_command = defaultdict(int)
        for time, entry in times:
            lexer = shlex.split(entry)
            times_by_command[lexer[0]] += time
        return self._top_n(times_by_command)

    def _slowest_commands(self, times, n=10):
        """执行时间较长top 10"""
        sorted_times = self._get_or_sort_list(times)
        slowest_commands = reversed(sorted_times[-n:])
        table_commands = [(str(time), self._reformat_entry(entry)) for time, entry in slowest_commands]
        return table_commands

    def _general_stats(self):
        """通用统计"""
        total_time = (self.last_ts - self.start_ts) / (1000*1000) or 1
        return [("Total Counts", self.line_count),
                ("Commands/Sec", '%.2f' % (self.line_count / total_time))]

    def _top_n(self, stat, n=10):
        if len(stat) > 0:
            sorted_items = sorted(stat.items(), key=lambda x: x[1], reverse=True)
            return sorted_items[:n]
        return stat

    def _pretty_out(self, result, title, percentages=False):
        """格式化"""
        data = "%s\r\n%s\r\n" % (title, '='*40)
        if not result:
            data += 'n/a\r\n'
            return data
        max_key_len = max((len(x[0]) for x in result))
        max_val_len = max((len(str(x[1])) for x in result))
        for key, val in result:
            key_padding = max(max_key_len-len(key), 0) * ' '
            if percentages:
                val_padding = max(max_val_len-len(str(val)), 0) * ' '
                val = '%s%s (%.2f%%)' % (val, val_padding, (float(val) / self.line_count) * 100)
            data += "%s%s\t%s\r\n" % (key, key_padding, val)
        return data

    def stats(self):
        """输出统计"""
        general = self._pretty_out(self._general_stats(), 'Overall Stats')
        prefixes = self._pretty_out(self._top_n(self.prefixes), 'Top Prefixes Keys', percentages=True)
        keys = self._pretty_out(self._top_n(self.keys), 'Top Keys', percentages=True)
        commands = self._pretty_out(self._top_n(self.commands), 'Top Commands', percentages=True)
        time_stats = self._pretty_out(self._time_stats(self.times), 'Command Time (microsecs)')
        heaviest_commands = self._pretty_out(self._heaviest_commands(self.times), 'Heaviest Commands (microsecs)')
        slowest_commands = self._pretty_out(self._slowest_commands(self.times), 'Slowest Calls')
        str = "%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s" % (general, prefixes, keys, commands, time_stats, heaviest_commands, slowest_commands)
        return str

    def handle_data(self, data):
        """对数据进行处理"""
        if isinstance(data, list):
            self.line_count = len(data)
            for item in data:
                self._record_duration(item)
                lexer = shlex.split(item["command"])
                self._record_command(lexer[0])
                if len(lexer) > 1:
                    self._record_key(lexer[1])

        return self.stats()