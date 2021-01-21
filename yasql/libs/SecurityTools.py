# -*— coding: utf-8 -*-
# __author__ : pandonglin
from Crypto.Cipher import AES
from binascii import b2a_hex, a2b_hex


class KMS:
    def __init__(self):
        self.KEY = 'mB4mB0rQ4tA3aB4e'.encode('utf-8')  # 16 位数
        self.IV = 'hB2aA1cQ3fE2qO0f'.encode('utf-8')

    # 如果text不足16位的倍数就用空格补足为16位
    def add_to_16(self, text):
        if len(text.encode('utf-8')) % 16:
            add = 16 - (len(text.encode('utf-8')) % 16)
        else:
            add = 0
        text = text + ('\0' * add)
        return text.encode('utf-8')

    # 加密函数
    def encrypt(self, text):
        text = self.add_to_16(text)
        cryptos = AES.new(self.KEY, AES.MODE_CBC, self.IV)
        cipher_text = cryptos.encrypt(text)
        return b2a_hex(cipher_text).decode('utf-8')

    # 解密后，去掉补足的空格用strip() 去掉
    def decrypt(self, text):
        cryptos = AES.new(self.KEY, AES.MODE_CBC, self.IV)
        plain_text = cryptos.decrypt(a2b_hex(text))
        return bytes.decode(plain_text).rstrip('\0')