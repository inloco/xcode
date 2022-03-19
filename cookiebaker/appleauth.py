import requests


class AppleAuth:
    def __init__(self, widget_key, base_url='https://idmsa.apple.com/appleauth', session=requests.Session()):
        self.base_url = base_url
        self.session = session

        self.session.headers.update({
            'Accept': 'application/json',
            'X-Apple-Widget-Key': widget_key,
        })

    def post_auth_signin(self, account_name, password):
        body = {
            'accountName': account_name,
            'password': password,
            'rememberMe': True,
        }

        r = self.session.post(f'{self.base_url}/auth/signin', json=body)
        if r.status_code != 409:
            r.raise_for_status()

        self.session.headers.update({
            'X-Apple-ID-Session-Id': r.headers['X-Apple-ID-Session-Id'],
            'scnt': r.headers['scnt'],
        })

        return r.json()

    def get_auth(self):
        r = self.session.get(f'{self.base_url}/auth')
        r.raise_for_status()

        return r.json()

    def put_auth_verify_phone(self, mode, phone_number_id):
        body = {
            'mode': mode,
            'phoneNumber': {
                'id': phone_number_id,
            },
        }

        r = self.session.put(f'{self.base_url}/auth/verify/phone', json=body)
        r.raise_for_status()

        return r.json()

    def post_auth_verify_phone_securitycode(self, mode, phone_number_id, security_code):
        body = {
            'mode': mode,
            'phoneNumber': {
                'id': phone_number_id,
            },
            'securityCode': {
                'code': security_code,
            },
        }

        r = self.session.post(f'{self.base_url}/auth/verify/phone/securitycode', json=body)
        r.raise_for_status()

        return r.json()

    def get_auth_2sv_trust(self):
        r = self.session.get(f'{self.base_url}/auth/2sv/trust')
        r.raise_for_status()
