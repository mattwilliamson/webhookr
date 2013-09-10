import random
import time

from urlhash import URLHash


hasher = URLHash("23456789abcdefghjklmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")


def generate():
    """Make a short alphanumeric ID"""
    global hasher
    return hasher.encode(random.random() * time.time())
