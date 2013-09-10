import math


class URLHash(object):
    codeset  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    base = 62 # len(codeset)

    def __init__(self, codeset=None):
        if codeset:
            self.codeset = codeset
            self.base = len(codeset)

    def encode(self, id):
        hash = ""
        while id > 0:
            hash = self.codeset[int(id % self.base)] + hash
            id = math.floor(id / self.base)
        return hash

    def decode(self, encoded):
        id = 0
        for index, char in enumerate(encoded[::-1]):
            n = self.codeset.find(char)
            if n == -1:
                return 0 # Invalid hash
            id += n * math.pow(self.base, index)
        return int(id)

