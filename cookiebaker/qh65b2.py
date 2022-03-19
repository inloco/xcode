import requests


class QH65B2:
    def __init__(self, base_url='https://developer.apple.com/services-account/QH65B2', session=requests.Session()):
        self.base_url = base_url
        self.session = session

    def post_downloadws_listdownloads(self):
        r = self.session.post(f'{self.base_url}/downloadws/listDownloads.action')
        r.raise_for_status()

        return r.json()
