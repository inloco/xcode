import os


class Config:
    appleid_user = os.environ['APPLEID_USER']
    appleid_pass = os.environ['APPLEID_PASS']

    messagebird_accesskey = os.environ['MESSAGEBIRD_ACCESSKEY']
    messagebird_caller = os.environ['MESSAGEBIRD_CALLER']
    messagebird_callee = os.environ['MESSAGEBIRD_CALLEE']

    awss3_bucketname = os.environ['AWSS3_BUCKETNAME']
    awss3_objectkey = os.getenv('AWSS3_OBJECTKEY')
