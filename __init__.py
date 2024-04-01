import os

from pocketbase import Client
from pocketbase.client import FileUpload

username = os.environ.get("DESCRIBER_USERNAME")
password = os.environ.get("DESCRIBER_PASSWORD")
url = os.environ.get("DESCRIBER_URL")

client = Client(url)
client.collection("users").auth_with_password(username, password)
images = client.collection("images")

image = images.create(
    {"file": FileUpload(("image.jpg", open("image.jpg", "rb"))), })
for followupId in image.followups:
    followup = client.collection("followups").get_one(followupId)
    print(followup.text)

followups = client.collection("followups")
followup = followups.create(
    {"text": "What types of trees are in the image?"})
print("Sending followup question")
image = images.update(image.id, {"followups+": followup.id})
followup = followups.get_one(image.followups[-1])
print(followup.text)
