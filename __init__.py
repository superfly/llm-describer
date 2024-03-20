import os

from pocketbase import Client
from pocketbase.client import FileUpload

username = os.environ.get("DESCRIBER_USERNAME")
password = os.environ.get("DESCRIBER_PASSWORD")
url = os.environ.get("DESCRIBER_URL")

client = Client(url)
client.collection("users").auth_with_password(username, password)
image = client.collection("images").create(
    {"file": FileUpload(("image.jpg", open("image.jpg", "rb"))), })
for followupId in image.followups:
    followup = client.collection("followups").get_one(followupId)
    print(followup.text)
