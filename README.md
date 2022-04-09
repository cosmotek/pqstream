# pqstream
Postgres-backed Event Library for Go

## TODO
- Improve docs (about considerations, setup example, if pqstream works across sessions etc)
- Better logging, possibly with an optional Zerolog adapter
- Authorization hooks for subscription which provide access to the http req,
  res, and the topic requested in order to allow or block access to a topic list
- Better configurability? Different encoders (Msgpack, Protobuf etc)