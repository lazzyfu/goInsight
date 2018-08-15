# -*- coding:utf-8 -*-
# edit by fuzongfei
# 使用
# python test4.py -uincep_user -pFuzongfei_1991 -H172.17.101.40 -P3306 -dtest -D./ -egbk -txlsx -f aaa.sql
# python test4.py -uincep_user -pFuzongfei_1991 -H172.17.101.40 -P3306 -dtest -D./ -eutf-8 -ttxt -f aaa.sql

import datetime
import re
import os
import time
import logging
import argparse
import zipfile

import pymysql
import sqlparse
from openpyxl import Workbook

# Set logger
LOG_FORMAT = "%(asctime)s - %(levelname)s - %(message)s"
DATE_FORMAT = "%Y-%m-%d %H:%M:%S"
logging.basicConfig(level=logging.DEBUG, format=LOG_FORMAT, datefmt=DATE_FORMAT)
logger = logging.getLogger(__name__)


def main():
    parser = argparse.ArgumentParser(description='数据导出帮助信息, 支持海量数据的导出.')
    parser.add_argument('-u', '--user', required=True, type=str, help='远程数据库用户名')
    parser.add_argument('-p', '--password', required=True, type=str, help='远程数据库密码')
    parser.add_argument('-H', '--host', required=True, type=str, help='远程数据库主机名')
    parser.add_argument('-P', '--port', required=True, type=int, help='远程数据库端口')
    parser.add_argument('-d', '--database', required=True, type=str, help='指定导出数据所在的的库')
    parser.add_argument('-e', '--encoding', type=str, help='导出文件编码，默认为：utf-8，可选编码为：utf-8、gbk')
    parser.add_argument('-t', '--type', type=str, help='导出文件的格式，默认为：txt，可选格式为：txt、xlsx')
    parser.add_argument('-D', '--directory', required=True, type=str, help='导出文件存放的目录')
    parser.add_argument('-f', '--file', type=str, required=True, help='加载导出SQL的文件，支持N条SQL批量自动导出')
    args = parser.parse_args()

    logging.info('### 提示 ###')
    logging.info(f'用户名：{args.user}')
    logging.info(f'密码：{args.password}')
    logging.info(f'主机：{args.host}')
    logging.info(f'端口：{args.port}')
    logging.info(f'库：{args.database}')
    logging.info(f"编码：{args.encoding if args.encoding else 'utf-8'}")
    logging.info(f"格式：{args.type if args.type else 'txt'}")
    logging.info(f'读取SQL列表文件：{args.file}')

    export_data = ExportData(args)
    export_data.run()


class ExportData(object):
    def __init__(self, args):
        self.user = args.user
        self.password = args.password
        self.host = args.host
        self.port = int(args.port)
        self.database = args.database
        self.file = args.file
        self.encoding = args.encoding if args.encoding else 'utf-8'
        self.type = args.type if args.type else 'txt'
        self.directory = args.directory
        self.date_desc = datetime.datetime.now().strftime("%Y%m%d%H%M")
        self.conn = pymysql.connect(host=self.host,
                                    user=self.user,
                                    password=self.password,
                                    port=self.port,
                                    db=self.database,
                                    charset='utf8')

    def read_sql(self):
        # 读取sql文件，并提取sql，返回sql列表
        sql_list = []
        with open(self.file, 'r', encoding='utf-8') as file:
            file_parse = sqlparse.parse(file.read().strip())
            for sql in file_parse:
                # 使用strip移除sql尾部的;
                sql_list.append(sql.value.strip(';'))
        return sql_list

    def get_count(self, sql):
        # 查询当前SQL的返回的查询数量，返回分页的SQL列表
        # 分页数量：10000
        rewrite_rule = re.compile('([\s\S]*)SELECT([\s\S]*) FROM ([\s\S]*)', re.I)
        count_query = rewrite_rule.sub(r'SELECT count(*) as count FROM \3', sql)

        self.conn.cursorclass = pymysql.cursors.DictCursor
        with self.conn.cursor() as cursor:
            cursor.execute(count_query)
            count = cursor.fetchone()

        num = int(count['count'] / 10000) + 1
        logging.info('\n')
        logging.info(f'正在导出SQL：{sql.strip()}')
        logging.info('SQL导出记录总数：%d' % count['count'])
        logging.info('SQL分页数量(page rows：10000)：%d' % num)
        offset_list = [i * 10000 for i in range(num)]
        page_list = [' '.join((sql, f'LIMIT 10000 OFFSET {i}')) for i in offset_list]
        return page_list

    def compress_file(self, file):
        # 压缩文件
        dst_filename = os.path.basename(file) + '.zip'
        dst_file = os.path.join(self.directory, dst_filename)
        logging.info(f'正在压缩文件：{file} ---> {dst_file}')
        with zipfile.ZipFile(dst_file, 'w', allowZip64=True, compression=zipfile.ZIP_DEFLATED) as filezip:
            filezip.write(file)
        logging.info(f'删除源文件：{file}')
        os.remove(file) if os.path.exists(file) else None

    def export_txt(self, pages_list, num):
        # 导出成txt格式
        # num：保存文件的结尾_num标识，为str类型
        num = str(num)
        file = os.path.join(self.directory, f'result_{self.date_desc}_{num}.txt')
        with open(file, 'a', encoding=f'{self.encoding}') as f:
            self.conn.cursorclass = pymysql.cursors.Cursor
            with self.conn.cursor() as cursor:
                for sql in pages_list:
                    progress = '/'.join([str(pages_list.index(sql) + 1), str(len(pages_list))])
                    logging.info(f'分页导出进度：{progress}')
                    cursor.execute(sql)
                    for line in cursor.fetchall():
                        data = '|&|'.join([str(k) for k in line])
                        f.write(data + '\n')
                    time.sleep(0.1)
        self.compress_file(file)

    def export_xlsx(self, pages_list, num):
        # 导出成xlsx格式
        # num：保存文件的结尾_num标识，为str类型
        num = str(num)
        file = os.path.join(self.directory, f'result_{self.date_desc}_{num}.xlsx')

        wb = Workbook()
        wb.encoding = f'{self.encoding}'
        ws = wb.active
        ws.title = f'result_{num}'

        # 获取列名作为标题
        self.conn.cursorclass = pymysql.cursors.DictCursor
        with self.conn.cursor() as cursor:
            cursor.execute(pages_list[0])
            title = []
            for column_name in cursor.fetchone():
                title.append(column_name)
        ws.append(title)

        # 获取数据，并写入到表格
        self.conn.cursorclass = pymysql.cursors.Cursor
        with self.conn.cursor() as cursor:
            for sql in pages_list:
                progress = '/'.join([str(pages_list.index(sql) + 1), str(len(pages_list))])
                logging.info(f'分页导出进度：{progress}')
                cursor.execute(sql)
                for row in cursor.fetchall():
                    ws.append(row)

        wb.save(file)
        self.compress_file(file)

    def run(self):
        sql_list = self.read_sql()
        for sql in sql_list:
            page_list = self.get_count(sql)
            start_time = time.time()
            if self.type == 'txt':
                self.export_txt(page_list, sql_list.index(sql))
                end_time = time.time()
            elif self.type == 'xlsx':
                self.export_xlsx(page_list, sql_list.index(sql))
                end_time = time.time()

            consume_time = ''.join((str(round(end_time - start_time, 2)), 's'))
            logging.info(f'耗时：{consume_time}')


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt as err:
        logging.warning('任务终止...')
    except Exception as err:
        logging.error(str(err))
