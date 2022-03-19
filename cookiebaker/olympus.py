import requests


class Olympus:
    def __init__(self, base_url='https://appstoreconnect.apple.com/olympus', session=requests.Session()):
        self.base_url = base_url
        self.session = session

    def get_app_config(self, hostname='itunesconnect.apple.com'):
        r = self.session.get(f'{self.base_url}/v1/app/config?hostname={hostname}')
        r.raise_for_status()

        return r.json()

    def get_session(self):
        r = self.session.get(f'{self.base_url}/v1/session')
        r.raise_for_status()

        return r.json()
