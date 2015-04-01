# broker
An experimental broker for the Radiodan ecosystem

## Starting

Compile binary, then run with `./broker --help` to see all available options.

## Protocol (WIP)
We presume (but unless specificed, do not enforce) JSON as the encoding format
for the `payload`.

### PubSub

By default, subscribers connect to port `7172` to recieve data. Publishers
connect to port `7173`.

#### Topic Subscriptions

Subscribers can recieve message for multiple topics. Topic names are in the
format `x.y.z`, where `.` is the delimiter. There can be any number of delimited
fields.

There is wildcard support for topic subscription, where `*` matches any topic
within the delimited field, and `#` matches any number of delimited fields.
These fields can be mixed in any order.

For example:

* `player.*.volume` will match any player's volume events
* `player.%` will match all player events

#### Sending a message to be published
    [ "TOPIC", "PAYLOAD" ]

#### Receiving a message from a subscription
    [ "SUBSCRIBED_TOPIC", "TOPIC_FROM_PUBLISHER", "PAYLOAD" ]

### Request / Reply

The default port is `7171`.
