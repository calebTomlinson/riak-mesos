package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/basho-labs/riak-mesos/process_manager"
	"github.com/basho-labs/riak-mesos/common"

	"bytes"
	"errors"
)

type DirectorNode struct {
	finishChan chan interface{}
	pm         *process_manager.ProcessManager
	running    bool
}

func NewDirectorNode() *DirectorNode {
	decompress()
	return &DirectorNode{
		running: false,
	}
}

func (directorNode *DirectorNode) runLoop() {
	waitChan := directorNode.pm.Listen()
	select {
	case <-waitChan:
		{
			log.Error("Director Died, failing")
		}
	case <-directorNode.finishChan:
		{
			log.Info("Finish channel says to shut down Director")
			directorNode.pm.TearDown()
		}
	}
	time.Sleep(15 * time.Second)
	log.Info("Shutting down")
}

func decompress() {
	var err error
	if err := os.Mkdir("director", 0777); err != nil {
		log.Fatal("Unable to make director directory: ", err)
	}

	asset, err := Asset("trusty.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	if err = common.ExtractGZ("director", bytes.NewReader(asset)); err != nil {
		log.Fatal("Unable to extract trusty root: ", err)
	}
	asset, err = Asset("riak_mesos_director-bin.tar.gz")

	if err != nil {
		log.Fatal(err)
	}
	if err = common.ExtractGZ("director", bytes.NewReader(asset)); err != nil {
		log.Fatal("Unable to extract rex: ", err)
	}
}

func (directorNode *DirectorNode) Run() {
	exepath := "/director/bin/director"

	var err error

	args := []string{"console", "-noinput"}
	healthCheckFun := func() error {
		log.Info("Checking is Director is started")
		logPath := filepath.Join(".", "director", "director", "log", "console.log")
		data, err := ioutil.ReadFile(logPath)
		if err != nil {
			if bytes.Contains(data, []byte("lager started on node")) {
				log.Info("Director started")
				return nil
			} else {
				return errors.New("Director not yet started")
			}
		} else {
			return err
		}
	}
	tearDownFun := func() {
		log.Info("Tearing down director")
	}

	log.Debugf("Starting up Director %v", exepath)

	chroot := filepath.Join(".", "director")
	directorNode.pm, err = process_manager.NewProcessManager(tearDownFun, exepath, args, healthCheckFun, &chroot)

	if err != nil {
		log.Error("Could not start Director: ", err)
		time.Sleep(15 * time.Minute)
		log.Info("Shutting down due to GC, after failing to bring up Director node")
	} else {
		directorNode.running = true
		directorNode.runLoop()
	}
}
