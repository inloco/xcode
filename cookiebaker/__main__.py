import requests
import sys

from .authenticator import *
from .cookiejar import *


def main():
    kargs = {}
    try:
        kargs['cookie_path'] = sys.argv[1]
    except IndexError:
        pass

    with CookieJar(**kargs) as cookie_jar:
        s = requests.Session()
        s.cookies = cookie_jar

        Authenticator.authenticate(s)


if __name__ == '__main__':
    main()
