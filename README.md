# LLM Describer

This describer connects with a remote [Ollama](https://ollama.com) server, sends
it images, and returns descriptions to users. A
[PocketBase](https://pocketbase.io) frontend provides authentication,
authorization, an admin interface, and other features you might expect when
building a SaaS frontend for AI models. Configurations for deploying to
[Fly.io](https://fly.io) are also included.

## Running

First you'll need an Ollama instance. The [llm/](llm/) directory contains a
simple configuration for deploying to Fly.io, but any instance will do. The
describer will run without an instance but all responses will be stubbed.

Next, set the `OLLAMA_API` environment variable to the URL of your instance.
With that done, run:

```bash
go run .
```

## Client usage

[__init__.py](__init__.py) is an example of how to use this describer from
Python. Any PocketBase client API should do, or you can use the built-in REST
APIs directly from any language.

First, create a user from your [admin interface](http://localhost:8090/_/). Then
set the following environment variables:

* `LLM_DESCRIBER_URL`
* `LLM_DESCRIBER_USER`
* `LLM_DESCRIBER_PASSWORD`

Download an image and place it in a file called *image.jpg* in the current
directory. Run the client, wait for the description, then type either a followup
question or "quit" to exit.
