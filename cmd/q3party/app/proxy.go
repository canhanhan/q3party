package app

import (
	"context"
	"time"

	"github.com/finarfin/q3party/pkg/masterclient"
	"github.com/finarfin/q3party/pkg/masterserver"
	"github.com/finarfin/q3party/pkg/protocol"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Starts a master server to proxy upstream servers",
	Run:   runProxy,
}

func init() {
	proxyCmd.PersistentFlags().StringP("bind", "b", "127.0.0.1:27950", "Bind address")
	proxyCmd.PersistentFlags().StringSlice("servers", nil, "Comma seperated list of master servers")
	viper.BindPFlag("bind", proxyCmd.PersistentFlags().Lookup("bind"))
	viper.BindPFlag("servers", proxyCmd.PersistentFlags().Lookup("servers"))
}

type MasterServerReader struct {
	ctx    context.Context
	cancel context.CancelFunc
	mc     *masterclient.MasterClient
	ch     chan []string
}

func NewMasterServerReader(ctx context.Context, servers []string) (*MasterServerReader, error) {
	ch := make(chan []string)
	ctx, cancel := context.WithCancel(ctx)

	mc, err := masterclient.NewMasterClient(ctx, servers, ch)
	if err != nil {
		return nil, err
	}

	return &MasterServerReader{
		ctx:    ctx,
		cancel: cancel,
		mc:     mc,
		ch:     ch,
	}, nil
}

func (msr *MasterServerReader) Close() error {
	msr.cancel()
	return msr.mc.Close()
}

func (msr *MasterServerReader) Servers(req *protocol.GetServersRequest) (*protocol.GetServersResponse, error) {
	msr.mc.Refresh(req)

	ctx, _ := context.WithDeadline(msr.ctx, time.Now().Add(2*time.Second))
	readCh := make(chan []string)
	servers := []string{}
	go func(readCh chan<- []string) {
		for {
			select {
			case v := <-msr.ch:
				readCh <- v
			case <-ctx.Done():
				return
			}
		}
	}(readCh)

	for {
		select {
		case v := <-readCh:
			if len(v) == 0 {
				break
			}

			for _, s := range v {
				servers = append(servers, s)
			}
		case <-msr.ctx.Done():
			break
		}
	}

	return &protocol.GetServersResponse{Servers: servers}, nil
}

func runProxy(cmd *cobra.Command, args []string) {
	log.Trace("Starting proxy command")

	bind := viper.GetString("bind")
	servers := viper.GetStringSlice("servers")
	msr, err := NewMasterServerReader(cmd.Context(), servers)
	if err != nil {
		log.Error(err)
		cmd.PrintErr(err)
		return
	}

	ms, err := masterserver.NewMasterServer(cmd.Context(), bind, msr)
	if err != nil {
		log.Error(err)
		cmd.PrintErr(err)
		return
	}
	defer ms.Close()

	<-cmd.Context().Done()
}
