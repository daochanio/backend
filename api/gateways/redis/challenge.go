package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/daochanio/backend/api/entities"
)

const ChallengeNamespace = "challenge"

func (r *redisGateway) GetChallengeByAddress(ctx context.Context, address string) (entities.Challenge, error) {
	challengeStr, err := r.client.Get(ctx, getFullKey(ChallengeNamespace, address)).Result()

	if err != nil {
		return entities.Challenge{}, fmt.Errorf("get challenge %w", err)
	}

	chalJson := &challengeJson{}
	err = json.Unmarshal([]byte(challengeStr), chalJson)

	if err != nil {
		return entities.Challenge{}, fmt.Errorf("decode challenge %w", err)
	}

	challenge := fromChallengeJson(chalJson)

	return challenge, nil
}

func (r *redisGateway) SaveChallenge(ctx context.Context, challenge entities.Challenge) error {
	challengeJson := toChallengeJson(challenge)
	bytes, err := json.Marshal(challengeJson)

	if err != nil {
		return fmt.Errorf("encode challenge %w", err)
	}

	_, err = r.client.Set(ctx, getFullKey(ChallengeNamespace, challenge.Address()), bytes, challenge.TTL()).Result()

	if err != nil {
		return fmt.Errorf("save challenge %w", err)
	}

	return nil
}

type challengeJson struct {
	Address string    `json:"address"`
	Message string    `json:"message"`
	Expires time.Time `json:"expires"`
}

func toChallengeJson(challenge entities.Challenge) challengeJson {
	return challengeJson{
		Address: challenge.Address(),
		Message: challenge.Message(),
		Expires: challenge.Expires(),
	}
}

func fromChallengeJson(challengeJson *challengeJson) entities.Challenge {
	return entities.NewChallenge(challengeJson.Address, challengeJson.Message, challengeJson.Expires)
}
