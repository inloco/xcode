import functools
import re
import requests.exceptions

from .appleauth import *
from .aws import *
from .config import *
from .olympus import *
from .qh65b2 import *
from .messagebird import *


class Authenticator:
    __code_re = re.compile(r'((?:\d\W*?){6})')
    __notd_re = re.compile(r'\D')

    @staticmethod
    def authenticate(session):
        olympus = Olympus(session=session)
        try:
            olympus.get_session()
        except requests.exceptions.HTTPError as e:
            if e.response.status_code != 401:
                raise e

            app_config = olympus.get_app_config()
            auth_service_url = app_config['authServiceUrl']
            auth_service_key = app_config['authServiceKey']

            Authenticator.__authenticate(olympus.session, auth_service_url, auth_service_key)
            olympus.get_session()

        qh65b2 = QH65B2(session=session)
        qh65b2.post_downloadws_listdownloads()

    @staticmethod
    def __authenticate(session, auth_service_url, auth_service_key):
        apple_auth = AppleAuth(auth_service_key, auth_service_url, session)

        try:
            Authenticator.__sfa(apple_auth)
        except requests.exceptions.HTTPError as e:
            if e.response.status_code != 409:
                raise e

            Authenticator.__mfa(apple_auth)

    @staticmethod
    def __sfa(apple_auth):
        apple_auth.post_auth_signin(Config.appleid_user, Config.appleid_pass)

    @staticmethod
    def __mfa(apple_auth):
        push_mode, phone_number_id = Authenticator.__request_call(apple_auth)

        time.sleep(15)

        recording_name, recording_content = Authenticator.__get_recording()

        transcript = Authenticator.__transcribe_recording(recording_name, recording_content)

        security_code = Authenticator.__extract_code(transcript)

        Authenticator.__validate_code(apple_auth, push_mode, phone_number_id, security_code)

    @staticmethod
    def __request_call(apple_auth):
        auth = apple_auth.get_auth()
        trusted_phone_numbers = auth['trustedPhoneNumbers']
        trusted_phone_numbers = filter(lambda e: e['pushMode'] == 'voice', trusted_phone_numbers)
        trusted_phone_number = next(trusted_phone_numbers)
        trusted_phone_number_id = trusted_phone_number['id']
        trusted_phone_number_push_mode = trusted_phone_number['pushMode']

        apple_auth.put_auth_verify_phone(trusted_phone_number_push_mode, trusted_phone_number_id)

        return trusted_phone_number_push_mode, trusted_phone_number_id

    @staticmethod
    def __get_recording():
        messagebird = MessageBird(Config.messagebird_accesskey)

        while True:
            calls = messagebird.get_calls()
            calls = filter(lambda e: e['source'] == Config.messagebird_caller and e[
                'destination'] == Config.messagebird_callee, calls)
            call = next(calls)
            call_id = call['id']
            call_status = call['status']

            if call_status != 'ended':
                continue

            legs = messagebird.get_legs(call_id)
            legs = filter(lambda e: e['source'] == Config.messagebird_callee, legs)
            leg = next(legs)
            leg_id = leg['id']

            recordings = messagebird.get_recordings(call_id, leg_id)
            recording = recordings[0]
            recording_id = recording['id']
            recording_format = recording['format']
            recording_status = recording['status']

            if recording_status != 'done':
                continue

            audio = messagebird.get_audio(recording_id, recording_format)

            return f'{recording_id}.{recording_format}', audio

    @staticmethod
    def __transcribe_recording(recording_name, recording_content):
        aws = AWS()

        object_key = Config.awss3_objectkey or recording_name

        aws.upload_file(Config.awss3_bucketname, object_key, recording_content)
        transcript = aws.transcribe_audio(Config.awss3_bucketname, object_key)

        return transcript

    @staticmethod
    def __extract_code(transcript):
        security_codes = Authenticator.__code_re.findall(transcript)
        security_codes = map(lambda e: Authenticator.__notd_re.sub('', e), security_codes)
        security_code = functools.reduce(lambda prev, curr: curr if not prev or prev == curr else None, security_codes)

        return security_code

    @staticmethod
    def __validate_code(apple_auth, push_mode, phone_number_id, security_code):
        apple_auth.post_auth_verify_phone_securitycode(push_mode, phone_number_id, security_code)
        apple_auth.get_auth_2sv_trust()
