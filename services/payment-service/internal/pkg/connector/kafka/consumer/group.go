package consumer

import "github.com/IBM/sarama"

type groupSubscriber struct {
	messageHandler MessageHandler
	setup          func(session sarama.ConsumerGroupSession) error
	cleanup        func(session sarama.ConsumerGroupSession) error
}

func (g groupSubscriber) Setup(session sarama.ConsumerGroupSession) error {
	if g.setup != nil {
		return g.setup(session)
	}
	return nil
}

func (g groupSubscriber) Cleanup(session sarama.ConsumerGroupSession) error {
	if g.cleanup != nil {
		return g.cleanup(session)
	}
	return nil
}

func (g groupSubscriber) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			err := g.messageHandler(ctx, session, message)
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if err != nil {
				return err
			}

			session.MarkMessage(message, "")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
