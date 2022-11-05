package docker

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"io"
	"math/big"

	"github.com/brist-ai/monoid/monoidprotocol"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/rs/zerolog/log"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[j.Int64()]
	}
	return string(b)
}

func ContainerWait(ctx context.Context, cli *client.Client, containerID string) (chan int, chan error) {
	done := make(chan int)
	errc := make(chan error)
	statusCh, errCh := cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	go func() {
		select {
		case err := <-errCh:
			if err != nil {
				errc <- err
				return
			}
		case status := <-statusCh:
			done <- int(status.StatusCode)
		}
	}()

	return done, errc
}

func ContainerLogs(ctx context.Context, cli *client.Client, containerID string, stdout bool, stderr bool) (io.Reader, error) {
	var buf bytes.Buffer

	res, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: stdout,
		ShowStderr: stderr,
	})
	if err != nil {
		return nil, err
	}

	stdcopy.StdCopy(&buf, &buf, res)
	defer res.Close()

	return &buf, nil
}

func readMessages(stream chan []byte, closer io.Closer) chan monoidprotocol.MonoidMessage {
	messageChan := make(chan monoidprotocol.MonoidMessage)
	go func() {
		for s := range stream {
			msg := monoidprotocol.MonoidMessage{}
			if err := json.Unmarshal(s, &msg); err != nil {
				log.Err(err).Msgf("Error unmarshalling message: %s", string(s))
				continue
			}

			messageChan <- msg
		}

		closer.Close()
		close(messageChan)
	}()

	return messageChan
}

func readRecords(stream chan monoidprotocol.MonoidMessage) chan monoidprotocol.MonoidRecord {
	recordChan := make(chan monoidprotocol.MonoidRecord)
	go func() {
		for s := range stream {
			if s.Type != monoidprotocol.MonoidMessageTypeRECORD || s.Record == nil {
				log.Debug().Msgf("Message type is not record: %s", string(s.Type))
			}

			recordChan <- *s.Record
		}

		close(recordChan)
	}()

	return recordChan

}
