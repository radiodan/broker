package pubsub

import (
	zmq "github.com/pebbe/zmq4"
	"log"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type Topic map[string]string

var wildTopicRegex *regexp.Regexp

func (b *Broker) Poll() {
	b.connect()
	runtime.SetFinalizer(b, (*Broker).Close)

	poller := zmq.NewPoller()
	poller.Add(b.PubSocket, zmq.POLLIN)
	poller.Add(b.SubSocket, zmq.POLLIN)

	// listen to all incoming messages
	b.SubSocket.SetSubscribe("")

	// empty set of topics for wildcard matching
	topics := make(Topic)

	for {
		polled, err := poller.Poll(1000 * time.Millisecond)

		if err != nil {
			log.Printf("!: %q", err)
			break
		}

		for _, item := range polled {
			switch socket := item.Socket; socket {
			case b.SubSocket:
				msg, _ := b.SubSocket.RecvMessage(0)
				topic := msg[0]
				data := msg[1]

				log.Printf("Topic: %s, Msg: %s", topic, data)
				b.PubSocket.SendMessage(append([]string{topic}, msg...))

				// iterate through wildcards, looking for matches
				for wildTopic, topicRegex := range topics {
					log.Printf("wildTopic: %s, topicRegex: %s", wildTopic, topicRegex)
					matched, _ := regexp.MatchString(topicRegex, topic)
					if matched == true {
						b.PubSocket.SendMessage([]string{wildTopic, topic, data})
					}
				}
			case b.PubSocket:
				msg, _ := b.PubSocket.RecvMessage(0)

				frame := msg[0]
				topic := frame[1:]

				switch frame[0] {
				case 1:
					if topicIsWild(topic) == true {
						topics = appendIfMissing(topics, topic)
					}
					log.Printf("Subscribed: %s", topic)
				case 0:
					log.Printf("UnSubscribed: %s", topic)
					if topicIsWild(topic) == true {
						delete(topics, topic)
					}
				}
			}
		}
	}
}

func topicIsWild(topic string) bool {
	// start of line, or a period,
	// followed by * or #
	// followed by end of line or period
	if wildTopicRegex == nil {
		wildTopicRegex, _ = regexp.Compile("(^|\\.){1}[\\*#]{1}($|\\.){1}")
	}

	matched := wildTopicRegex.MatchString(topic)
	return matched
}

func appendIfMissing(topics Topic, topic string) Topic {
	_, exists := topics[topic]

	// already been added, somehow
	if exists == true {
		return topics
	}

	topicRegex := []string{}

	for _, e := range strings.Split(topic, ".") {
		var quoted string

		switch e {
		case "*":
			quoted = "[^\\.]*"
		case "#":
			quoted = ".*"
		default:
			quoted = regexp.QuoteMeta(e)
		}

		topicRegex = append(topicRegex, quoted)
	}

	topics[topic] = "^" + strings.Join(topicRegex, "\\.") + "$"

	return topics
}
