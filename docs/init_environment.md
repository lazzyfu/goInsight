# 初始化环境
首先您需要按照下面步骤安装下系统依赖包，强烈建议您选择一个干净的系统。

#### 安装系统依赖包
```bash
yum -y install epel-release

yum -y install \
net-tools bzip2-devel gcc gcc-c++ make automake unzip curl curl-devel \
libffi-devel perl perl-devel expat expat-devel zlib zlib-devel asciidoc \
xmlto gettext-devel openssl openssl-devel mlocate python-devel openldap-devel \
readline-devel git mysql-devel p7zip
```

#### 安装python3.7
```bash
wget https://www.python.org/ftp/python/3.7.9/Python-3.7.9.tgz
tar -zxf Python-3.7.9.tgz
cd Python-3.7.9/
./configure
make -j 4 && make install
```

#### 创建python3.7虚拟环境
```bash
/usr/local/bin/python3.7 -m pip install --upgrade pip
/usr/local/bin/pip3.7 install virtualenv -i https://mirrors.aliyun.com/pypi/simple
/usr/local/bin/virtualenv /venvyasql --python=/usr/local/bin/python3.7
```