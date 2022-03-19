import requests


class MessageBird:
    def __init__(self, access_key, base_url='https://voice.messagebird.com', session=requests.Session()):
        self.base_url = base_url
        self.session = session

        self.session.headers.update({
            'Authorization': f'AccessKey {access_key}',
        })

    def get_calls(self):
        r = self.session.get(f'{self.base_url}/calls')
        r.raise_for_status()

        j = r.json()
        return j['data']

    def get_legs(self, call_id):
        r = self.session.get(f'{self.base_url}/calls/{call_id}/legs')
        r.raise_for_status()

        j = r.json()
        return j['data']

    def get_recordings(self, call_id, leg_id):
        r = self.session.get(f'{self.base_url}/calls/{call_id}/legs/{leg_id}/recordings')
        r.raise_for_status()

        j = r.json()
        return j['data']

    def get_audio(self, recording_id, recording_format):
        r = self.session.get(f'{self.base_url}/recordings/{recording_id}.{recording_format}')
        r.raise_for_status()

        return r.content
