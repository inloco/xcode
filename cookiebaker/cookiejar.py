import http.cookiejar


class CookieJar(http.cookiejar.MozillaCookieJar):
    def __init__(self, cookie_path='cookies.txt'):
        super(CookieJar, self).__init__(cookie_path)

    def __enter__(self):
        try:
            self.load()
        except FileNotFoundError:
            self.save()

        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        if exc_type or exc_val or exc_tb:
            return False

        self.save(ignore_discard=True, ignore_expires=True)
        return True
