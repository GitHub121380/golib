package mcpacknpc

import (
	"bytes"
	"context"

	"github.com/GitHub121380/golib/gomcpack/mcpack"
	"github.com/GitHub121380/golib/gomcpack/npc"
)

type Client struct {
	*npc.Client
}

func NewClient(server []string) *Client {
	c := npc.NewClient(server)
	return &Client{Client: c}
}

func (c *Client) Call(ctx context.Context, args interface{}, reply interface{}) error {
	content, err := mcpack.Marshal(args)
	if err != nil {
		return err
	}
	resp, err := c.Client.Do(npc.NewRequest(ctx, bytes.NewReader(content)))
	if err != nil {
		return err
	}
	return mcpack.Unmarshal(resp.Body, reply)
}

func (c *Client) Send(ctx context.Context, args []byte) ([]byte, error) {
	resp, err := c.Client.Do(npc.NewRequest(ctx, bytes.NewReader(args)))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
