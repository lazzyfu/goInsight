import socket
from AuditSQL import settings

inception_host = getattr(settings, 'INCEPTION_HOST')
inception_port = getattr(settings, 'INCEPTION_PORT')

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
result = sock.connect_ex((inception_host, inception_port))

if 0 == result:
    print("Port is open")
else:
    print("Port is not open，return code：%s" % result)