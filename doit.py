# third parties imports
from googleapiclient.discovery import build
from httplib2 import Http
from oauth2client.service_account import ServiceAccountCredentials
import pprint

# modify according to your requirements
CLIENT_SECRET = 'credentials.json'  # the credentials downloaded from the GCP Console
ADMIN_USER = 'soenke.ruempler@superluminar.io'  # The admin user used by the service account

# Scopes
SCOPES = [
    'https://www.googleapis.com/auth/admin.directory.user',
    'https://www.googleapis.com/auth/admin.directory.group',
]

# service account initialization
credentials = ServiceAccountCredentials.from_json_keyfile_name(CLIENT_SECRET, scopes=SCOPES)
delegated_admin = credentials.create_delegated(ADMIN_USER)
admin_http_auth = delegated_admin.authorize(Http())

admin_sdk = build('admin', 'directory_v1', http=admin_http_auth)  # Admin SDK service

gsuite_users = admin_sdk.users().list(domain='superluminar.io').execute()['users']

def get_groups_of_user(user):
    raw_groups = admin_sdk.groups().list(userKey=user).execute()
    return [group['name'] for group in raw_groups.get("groups", []) if group["name"].startswith("AWS")]

def get_custom_schema_for_user(user):
    result = []
    groups = get_groups_of_user(user)
    for g in groups:
        account, role = g.split("-")[1:]
        result.append({
            "type": "work",
            "value": "arn:aws:iam::{aws_account_id}:role/{role},arn:aws:iam::{aws_account_id}:saml-provider/Google".format(aws_account_id=account, role=role)
        })
    return result

for user in gsuite_users:
    email = user["primaryEmail"]
    user_patch_params = {
        'userKey': email,
        'body': {
            "customSchemas": {
                "SSO": {
                    "session_duration": "43200",
                    "role": get_custom_schema_for_user(email)
                }
            }
        }
    }

    user_patch = admin_sdk.users().patch(**user_patch_params).execute()
    pprint.pprint(admin_sdk.users().get(userKey=email, projection="full").execute())
