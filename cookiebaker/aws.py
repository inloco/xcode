import boto3
import io
import time
import requests


class AWS:
    def __init__(self, session=requests.Session()):
        self.session = session

        self.s3 = boto3.client('s3')
        self.transcribe = boto3.client('transcribe')

    def upload_file(self, bucket_name, object_key, content):
        stream = io.BytesIO(content)
        self.s3.upload_fileobj(stream, bucket_name, object_key)

    def transcribe_audio(self, bucket_name, object_key):
        transcription_job_name = AWS.__calc_transcription_job_name(object_key)

        self.transcribe.start_transcription_job(
            TranscriptionJobName=transcription_job_name,
            LanguageCode='en-US',
            Media={
                'MediaFileUri': f's3://{bucket_name}/{object_key}',
            },
        )

        while True:
            transcription_job = self.transcribe.get_transcription_job(
                TranscriptionJobName=transcription_job_name,
            )
            transcription_job = transcription_job['TranscriptionJob']

            transcription_job_status = transcription_job['TranscriptionJobStatus']
            if transcription_job_status == 'COMPLETED':
                break
            if transcription_job_status == 'FAILED':
                raise Exception(transcription_job['FailureReason'])

            time.sleep(1)

        transcript = transcription_job['Transcript']
        transcript_file_uri = transcript['TranscriptFileUri']

        r = self.session.get(transcript_file_uri)
        r.raise_for_status()

        j = r.json()
        results = j['results']
        transcripts = results['transcripts']
        transcript = transcripts[0]['transcript']

        return transcript

    @staticmethod
    def __calc_transcription_job_name(object_key):
        name = object_key
        name = name[name.rfind('/') + 1:]
        name = name[:(name.find('.') + 1 or len(name) + 1) - 1]
        return name
