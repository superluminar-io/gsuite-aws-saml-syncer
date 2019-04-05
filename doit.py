for user in gsuite_users
    WRITE NEW CUSTOM ATTRIBUTES
    # we set a value for our custom schema on the desired user
    groups = get_groups_of_user()
    user_patch_params = {
        'userKey': SCHEMA_USER,
        'body': {
            "customSchemas": {
  "SSO": {
   "session_duration": "43200",
   "role": [
    {
     "type": "work",
     "value": "arn:aws:iam::197726340368:role/GoogleAccess,arn:aws:iam::197726340368:saml-provider/Google"
    },
    {
     "type": "work",
     "value": "arn:aws:iam::208513435484:role/GoogleAccess,arn:aws:iam::208513435484:saml-provider/Google"
    },
    {
     "type": "work",
     "value": "arn:aws:iam::197726340368:role/ReadOnly,arn:aws:iam::197726340368:saml-provider/Google"
    },
    {
     "type": "work",
     "value": "arn:aws:iam::197726340368:role/A,arn:aws:iam::197726340368:saml-provider/Google"
    },
    {
     "type": "work",
     "value": "arn:aws:iam::197726340368:role/C,arn:aws:iam::197726340368:saml-provider/Google"
    },
    {
     "type": "work",
     "value": "arn:aws:iam::197726340368:role/D,arn:aws:iam::197726340368:saml-provider/Google"
    },
    {
     "type": "work",
     "value": "arn:aws:iam::197726340368:role/B,arn:aws:iam::197726340368:saml-provider/Google"
    },
    {
     "type": "work",
     "value": "arn:aws:iam::197726340368:role/E,arn:aws:iam::197726340368:saml-provider/Google"
    }
   ]
  }

        },
    }
    user_patch = admin_sdk.users().patch(**user_patch_params).execute()

