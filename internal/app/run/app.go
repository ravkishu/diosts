package run

import (
	"bufio"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hakluke/haksecuritytxt/pkg/securitytxt"
)

type App struct {
	*Config

	reader *Reader
	writer *Writer
	pool *WorkerPool
}

func New(config *Config) (*App, error) {
	domainCh := make(chan string, config.NumThreads)
	txtCh := make(chan *securitytxt.SecurityTxt, config.NumThreads)

	reader, err := NewReader(config, domainCh)
	if err != nil {
		return nil, err
	}

	writer, err := NewWriter(config, txtCh)
	if err != nil {
		return nil, err
	}

	client, err := securitytxt.NewDomainClient(config.SecurityTxt)
	if err != nil {
		return nil, err
	}

	pool, err := NewWorkerPool(config, client, domainCh, txtCh)
	if err != nil {
		return nil, err
	}

	a := &App{
		Config: config,
		reader: reader,
		writer: writer,
		pool: pool,
	}

	return a, nil
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	go func() {
		err <- errCh
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
	}

	// TODO
	// - Add counter: total input. total securitytxt, total valid
	// - validate input
	// - error output for non-existent domain etc
	// - graceful shutdown:w
	// - output to stdout or file
	// - silent mode or log to stderr

	// Read domains - will close domainCh when done
	err := a.reader.Start(errCh)
	if err != nil {
		return err
	}

	// Write domains - will stop when txtCh is closed
	err = a.writer.Start()
	if err != nil {
		return err
	}

	// Will run until domainCh is closed and closes txtCh
	return a.pool.Run(errCh)
}
