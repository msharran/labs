# Min required fields for each item type
"""
{
        "folders": [
            {
                "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
                "name": "Folder Name"
                },
            ],
        "items": [
            {
                "type": 1,
                "name": "Login Item's Name",
                "folderId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
                "fields": [
                    {
                        "name": "custom-field-1",
                        "value": "custom-field-value",
                        "type": 0
                        },
                    ],
                "login": {
                    "uris": [
                        {
                            "match": null,
                            "uri": "https://mail.google.com"
                            }
                        ],
                    "username": "myaccount@gmail.com",
                    "password": "myaccountpassword",
                    ]
                }
            },
        ]
"""

"""
for file in "~/.password-store/**/*.gpg":
    Example, file = ~/.password-store/Email/foo/bar/baz.gpg
    get pass by running `$ pass Email/foo/bar/baz`. Output is Multi-line string returned by pass command
    first directory in the path is folder. Example, ~/.password-store/Email/foo/bar/baz.gpg => Email
    first line is password at /items/n/login/password
    remaining lines are custom fields at /items/n/fields
    custom field with login is /items/n/login/username
    custom field with url is   /items/n/login/uris/0/uri
"""




import os
import subprocess
import json
def get_pass(file):
    return subprocess.check_output(["pass", file]).decode("utf-8")


def get_folders(file):
    return file.split("/")[2]


def get_name(file):
    return file.split("/")[-1].split(".")[0]


def get_fields(pass):
    return [{"name": f[0], "value": f[1]} for f in pass.split("\n")[1:]]
